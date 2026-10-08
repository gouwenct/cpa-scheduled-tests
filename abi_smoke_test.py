import base64, ctypes, json, os, sys, tempfile, time
from ctypes import *

SO = os.path.abspath(sys.argv[1]) if len(sys.argv) > 1 else os.path.join(os.path.dirname(__file__), 'dist', 'cpa-scheduled-tests.so')
lib = ctypes.CDLL(SO)
libc = ctypes.CDLL('msvcrt.dll') if os.name == 'nt' else ctypes.CDLL(None)
libc.malloc.restype = c_void_p
libc.malloc.argtypes = [c_size_t]
libc.free.argtypes = [c_void_p]

class Buffer(Structure):
    _fields_ = [('ptr', POINTER(c_uint8)), ('len', c_size_t)]

HOSTCALL = CFUNCTYPE(c_int, c_void_p, c_char_p, POINTER(c_uint8), c_size_t, POINTER(Buffer))
HOSTFREE = CFUNCTYPE(None, c_void_p, c_size_t)
PLUGCALL = CFUNCTYPE(c_int, c_char_p, POINTER(c_uint8), c_size_t, POINTER(Buffer))
PLUGFREE = CFUNCTYPE(None, c_void_p, c_size_t)
PLUGSHUT = CFUNCTYPE(None)

class HostAPI(Structure):
    _fields_ = [('abi_version', c_uint32), ('host_ctx', c_void_p), ('call', HOSTCALL), ('free_buffer', HOSTFREE)]
class PluginAPI(Structure):
    _fields_ = [('abi_version', c_uint32), ('call', PLUGCALL), ('free_buffer', PLUGFREE), ('shutdown', PLUGSHUT)]

accounts = [
    {'id':'codex-a','auth_index':'a1','name':'codex-a.json','provider':'codex','email':'a@example.com','status':'ready','disabled':False,'unavailable':False},
    {'id':'codex-b','auth_index':'b2','name':'codex-b.json','provider':'codex','email':'b@example.com','status':'disabled','disabled':True,'unavailable':True},
]
http_calls = []

def alloc_json(obj, out):
    b = json.dumps({'ok': True, 'result': obj}, separators=(',', ':')).encode()
    p = libc.malloc(len(b))
    ctypes.memmove(p, b, len(b))
    out.contents.ptr = cast(p, POINTER(c_uint8))
    out.contents.len = len(b)

@HOSTCALL
def host_call(ctx, method, req_ptr, req_len, out):
    method = method.decode()
    req = bytes(string_at(req_ptr, req_len)) if req_ptr and req_len else b'{}'
    payload = json.loads(req or b'{}')
    if method == 'host.auth.list':
        alloc_json({'files': accounts}, out)
        return 0
    if method == 'host.auth.get':
        idx = payload.get('auth_index')
        alloc_json({'auth_index':idx,'name':idx+'.json','json':{'type':'codex','access_token':'token-'+idx,'account_id':'acct-'+idx}}, out)
        return 0
    if method == 'host.http.do':
        http_calls.append(payload)
        if payload.get('Method') == 'GET':
            assert payload['URL'] == 'https://chatgpt.com/backend-api/wham/usage', payload['URL']
            body = {'rate_limit': {'primary_window': {
                'limit_window_seconds': 18000, 'used_percent': 25.0, 'reset_at': 4102444800,
            }}}
            alloc_json({'StatusCode':200,'Headers':{},'Body':base64.b64encode(json.dumps(body).encode()).decode()}, out)
            return 0
        assert payload['URL'] == 'https://chatgpt.com/backend-api/codex/responses'
        request_body = json.loads(base64.b64decode(payload['Body']))
        assert request_body['input'][0]['content'][0]['text'] == '你好', request_body
        alloc_json({'StatusCode':200,'Headers':{
            'x-codex-primary-used-percent':['25.0'],
            'x-codex-primary-window-minutes':['300'],
            'x-codex-primary-reset-at':['4102444800'],
        } if payload['Headers']['Authorization'] == ['Bearer token-a1'] else {},'Body':base64.b64encode(b'data: ok').decode()}, out)
        return 0
    if method == 'host.log':
        alloc_json({'logged': True}, out)
        return 0
    b = json.dumps({'ok':False,'error':{'code':'unknown','message':method}}).encode()
    p = libc.malloc(len(b)); ctypes.memmove(p,b,len(b)); out.contents.ptr=cast(p,POINTER(c_uint8)); out.contents.len=len(b)
    return 0

@HOSTFREE
def host_free(ptr, n):
    libc.free(ptr)

host = HostAPI(1, None, host_call, host_free)
plugin = PluginAPI()
lib.cliproxy_plugin_init.argtypes = [POINTER(HostAPI), POINTER(PluginAPI)]
lib.cliproxy_plugin_init.restype = c_int
rc = lib.cliproxy_plugin_init(byref(host), byref(plugin))
assert rc == 0, rc
assert plugin.abi_version == 1

def call(method, payload=None):
    raw = json.dumps(payload or {}, separators=(',',':')).encode()
    arr = (c_uint8 * len(raw)).from_buffer_copy(raw) if raw else None
    out = Buffer()
    rc = plugin.call(method.encode(), cast(arr, POINTER(c_uint8)) if arr else None, len(raw), byref(out))
    data = bytes(string_at(out.ptr, out.len)) if out.ptr and out.len else b''
    if out.ptr: plugin.free_buffer(cast(out.ptr, c_void_p), out.len)
    assert rc == 0, (method, rc, data)
    obj = json.loads(data)
    assert obj.get('ok') is True, (method, obj)
    return obj['result']

test_data = tempfile.TemporaryDirectory()
reg = call('plugin.register', {'schema_version': 6, 'plugin_dir': test_data.name})
assert reg['schema_version'] == 6
assert reg['metadata']['Name'] == 'CPA Scheduled Tests'
assert reg['metadata']['GitHubRepository']
mg = call('management.register')
assert any(r['Path'].endswith('/run-all') for r in mg['routes'])

def mgmt(method, path, body=None, query=None):
    req = {
        'Method': method,
        'Path': '/v0/management' + path,
        'Headers': {},
        'Query': query or {},
        'Body': base64.b64encode(json.dumps(body or {}).encode()).decode() if body is not None else '',
    }
    result = call('management.handle', req)
    body_raw = base64.b64decode(result.get('Body','')) if result.get('Body') else b''
    parsed = json.loads(body_raw) if body_raw else None
    return result['StatusCode'], parsed

status, st = mgmt('GET','/plugins/cpa-scheduled-tests/state')
assert status == 200
assert len(st['accounts']) == 2
status, bulk = mgmt('POST','/plugins/cpa-scheduled-tests/plans/bulk-create', {
    'name_prefix':'All accounts',
    'model':'gpt-test',
    'cron':'0 6,11,16,21 * * *',
    'timezone':'Asia/Shanghai',
    'prompt':'ping',
    'enabled':True,
})
assert status == 200, (status, bulk)
assert bulk['total_accounts'] == 2 and bulk['created'] == 2 and bulk['skipped'] == 0, bulk
status, bulk2 = mgmt('POST','/plugins/cpa-scheduled-tests/plans/bulk-create', {
    'name_prefix':'All accounts',
    'model':'gpt-test',
    'cron':'0 6,11,16,21 * * *',
    'timezone':'Asia/Shanghai',
    'prompt':'ping',
    'enabled':True,
})
assert status == 200 and bulk2['created'] == 0 and bulk2['skipped'] == 2, bulk2
status, st = mgmt('GET','/plugins/cpa-scheduled-tests/state')
assert len(st['plans']) == 2, st['plans']
status, accepted = mgmt('POST','/plugins/cpa-scheduled-tests/run-all', {'model':'gpt-test','prompt':''})
assert status == 202, (status, accepted)
for _ in range(80):
    time.sleep(0.05)
    status, st = mgmt('GET','/plugins/cpa-scheduled-tests/state')
    if st['jobs'] and not st['jobs'][0]['running']:
        break
else:
    raise AssertionError('bulk job did not finish')
assert st['jobs'][0]['total'] == 2
assert st['jobs'][0]['succeeded'] == 2
assert len(http_calls) == 3, http_calls
assert sum(x['Method'] == 'POST' for x in http_calls) == 2
assert sum(x['Method'] == 'GET' for x in http_calls) == 1
status, logs = mgmt('GET','/plugins/cpa-scheduled-tests/logs', query={'limit':['20']})
assert status == 200
assert len(logs['logs']) >= 2
seen = {x['auth_index'] for x in logs['logs'][:2]}
assert {'a1','b2'}.issubset(seen), seen
assert any(x.get('disabled') for x in logs['logs'] if x['auth_index']=='b2')
assert all(x.get('five_hour_status') == 'open' for x in logs['logs'][:2]), logs['logs'][:2]
assert all(x.get('five_hour_used_percent') == 25.0 for x in logs['logs'][:2]), logs['logs'][:2]
plugin.shutdown()
test_data.cleanup()
print('ABI smoke test passed: registration, state, bulk plan create+dedupe, run-all, disabled account inclusion, logs')
