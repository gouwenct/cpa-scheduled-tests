package main

func panelHTML() string {
	return `<!doctype html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>CPA Scheduled Tests</title>
<style>
:root{color-scheme:light dark;--bg:#f5f7fb;--card:#fff;--text:#172033;--muted:#6b7280;--line:#e5e7eb;--accent:#2563eb;--accent2:#1d4ed8;--danger:#dc2626;--ok:#059669;--warn:#d97706;--chip:#eef2ff;--shadow:0 8px 24px rgba(15,23,42,.07)}
@media(prefers-color-scheme:dark){:root{--bg:#0f172a;--card:#111827;--text:#e5e7eb;--muted:#9ca3af;--line:#263244;--accent:#60a5fa;--accent2:#93c5fd;--danger:#f87171;--ok:#34d399;--warn:#fbbf24;--chip:#1e293b;--shadow:none}}
*{box-sizing:border-box}body{margin:0;background:var(--bg);color:var(--text);font:14px/1.5 Inter,ui-sans-serif,system-ui,-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif}.wrap{max-width:1440px;margin:0 auto;padding:28px}.top{display:flex;justify-content:space-between;gap:20px;align-items:flex-start;margin-bottom:20px}.title h1{font-size:24px;margin:0 0 4px}.title p{margin:0;color:var(--muted)}.keybox{display:flex;gap:8px;align-items:center;flex-wrap:wrap}.keybox input{width:260px}.cards{display:grid;grid-template-columns:repeat(4,minmax(150px,1fr));gap:12px;margin:16px 0 22px}.card,.section{background:var(--card);border:1px solid var(--line);border-radius:12px;box-shadow:var(--shadow)}.card{padding:16px}.metric{font-size:24px;font-weight:700}.label{font-size:12px;color:var(--muted);margin-top:2px}.section{margin:16px 0;padding:18px}.sectionhead{display:flex;justify-content:space-between;gap:12px;align-items:center;margin-bottom:14px}.sectionhead h2{font-size:17px;margin:0}.actions{display:flex;gap:8px;align-items:center;flex-wrap:wrap}button,input,select,textarea{font:inherit}input,select,textarea{background:var(--card);color:var(--text);border:1px solid var(--line);border-radius:8px;padding:9px 10px;outline:none}input:focus,select:focus,textarea:focus{border-color:var(--accent)}button{border:1px solid var(--line);background:var(--card);color:var(--text);border-radius:8px;padding:9px 12px;cursor:pointer}button:hover{border-color:var(--accent)}button.primary{background:var(--accent);border-color:var(--accent);color:white}button.primary:hover{background:var(--accent2)}button.danger{color:var(--danger)}button:disabled{opacity:.5;cursor:not-allowed}.grid2{display:grid;grid-template-columns:1fr 1fr;gap:12px}.grid4{display:grid;grid-template-columns:1.2fr 1.2fr 1fr auto;gap:10px;align-items:end}.field label{display:block;color:var(--muted);font-size:12px;margin:0 0 5px}.field input,.field select,.field textarea{width:100%}.hint{font-size:12px;color:var(--muted);margin-top:6px}.tablewrap{overflow:auto;border:1px solid var(--line);border-radius:10px}table{width:100%;border-collapse:collapse;min-width:980px}th,td{padding:10px 11px;border-bottom:1px solid var(--line);text-align:left;vertical-align:middle}th{font-size:12px;color:var(--muted);font-weight:600;background:color-mix(in srgb,var(--card) 94%,var(--muted) 6%)}tr:last-child td{border-bottom:0}.chip{display:inline-flex;align-items:center;border-radius:999px;padding:2px 8px;font-size:12px;background:var(--chip)}.ok{color:var(--ok)}.bad{color:var(--danger)}.warn{color:var(--warn)}.muted{color:var(--muted)}.nowrap{white-space:nowrap}.ellipsis{max-width:360px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}.toggle{display:inline-flex;align-items:center;gap:6px}.small{padding:6px 8px;font-size:12px}dialog{width:min(680px,92vw);border:1px solid var(--line);border-radius:12px;background:var(--card);color:var(--text);padding:0;box-shadow:0 24px 80px rgba(0,0,0,.25)}dialog::backdrop{background:rgba(0,0,0,.4)}.modalhead,.modalbody,.modalfoot{padding:16px 18px}.modalhead{border-bottom:1px solid var(--line);font-size:17px;font-weight:700}.modalfoot{border-top:1px solid var(--line);display:flex;justify-content:flex-end;gap:8px}.modalbody{display:grid;gap:12px}.toast{position:fixed;right:20px;bottom:20px;max-width:420px;background:#111827;color:#fff;padding:12px 14px;border-radius:10px;box-shadow:0 10px 30px rgba(0,0,0,.25);display:none;z-index:10}.dangerbox{border:1px solid color-mix(in srgb,var(--danger) 40%,var(--line));background:color-mix(in srgb,var(--danger) 6%,var(--card));padding:12px;border-radius:10px}.tabs{display:flex;gap:8px;margin-bottom:12px}.jobs{display:flex;gap:8px;flex-wrap:wrap}.job{border:1px solid var(--line);border-radius:8px;padding:7px 9px;font-size:12px}.mono{font-family:ui-monospace,SFMono-Regular,Menlo,Consolas,monospace}@media(max-width:900px){.wrap{padding:16px}.top{display:block}.keybox{margin-top:12px}.cards{grid-template-columns:repeat(2,1fr)}.grid2,.grid4{grid-template-columns:1fr}.keybox input{width:100%}}
</style>
</head>
<body>
<div class="wrap">
  <div class="top">
    <div class="title"><h1>Scheduled Tests</h1><p>CPA 原生的账号 × 模型 × Cron 定时测试。支持单账号立即发送、全部账号强制发送和持久化日志。</p></div>
    <div class="keybox"><input id="mgmtKey" type="password" autocomplete="off" placeholder="CPA Management Key"><button onclick="saveKey()">保存本会话</button><button onclick="loadModels(true)">同步模型</button><button class="primary" onclick="loadAll()">刷新</button></div>
  </div>

  <div class="cards">
    <div class="card"><div id="mAccounts" class="metric">-</div><div class="label">Codex 账号</div></div>
    <div class="card"><div id="mPlans" class="metric">-</div><div class="label">定时计划</div></div>
    <div class="card"><div id="mEnabled" class="metric">-</div><div class="label">启用计划</div></div>
    <div class="card"><div id="mRunning" class="metric">-</div><div class="label">运行中任务</div></div>
  </div>

  <div class="section">
    <div class="sectionhead"><h2>快速发送</h2><div class="muted">不创建计划，直接对指定账号发送一次真实请求</div></div>
    <div class="grid4">
      <div class="field"><label>账号</label><select id="quickAccount"></select></div>
      <div class="field"><label>Model</label><input id="quickModel" list="modelSuggestions" placeholder="例如 gpt-5.6-luna"></div>
      <div class="field"><label>最小 Prompt</label><input id="quickPrompt" placeholder="ping"></div>
      <button class="primary" onclick="quickRun()">立即发送</button>
    </div>
  </div>

  <div class="section">
    <div class="sectionhead"><h2>定时测试计划</h2><div class="actions"><button onclick="openBulkPlan()">一键添加全部账号</button><button class="primary" onclick="openPlan()">+ 新建计划</button></div></div>
    <div class="tablewrap"><table><thead><tr><th>启用</th><th>名称</th><th>账号</th><th>Model</th><th>Cron</th><th>时区</th><th>下次执行</th><th>上次结果</th><th>操作</th></tr></thead><tbody id="planRows"></tbody></table></div>
    <div class="hint">Cron 为标准 5 段：分钟 小时 日 月 星期。例如 <span class="mono">0 6,11,16,21 * * *</span>。</div>
  </div>

  <div class="section dangerbox">
    <div class="sectionhead"><h2>全部账号立即发送</h2><div class="muted">强制遍历 CPA 中所有 Codex auth 记录，包括 disabled / unavailable；失败账号会写入日志。</div></div>
    <div class="grid4">
      <div class="field"><label>Model</label><input id="allModel" list="modelSuggestions"></div>
      <div class="field"><label>Prompt</label><input id="allPrompt" placeholder="ping"></div>
      <div class="field"><label>并发</label><input id="bulkConcurrency" type="number" min="1" max="32"></div>
      <button class="primary" onclick="runAllAccounts()">一键向全部账号发送</button>
    </div>
    <div class="hint">该操作不会检查 5h 窗口是否已经启动；点击后就是强制发送。适合你需要统一点燃窗口的场景。</div>
  </div>

  <div class="section">
    <div class="sectionhead"><h2>运行任务</h2><div class="actions"><button onclick="loadAll()">刷新状态</button></div></div>
    <div id="jobs" class="jobs"><span class="muted">暂无</span></div>
  </div>

  <div class="section">
    <div class="sectionhead"><h2>日志</h2><div class="actions"><select id="logLimit" onchange="loadLogs()"><option>100</option><option selected>200</option><option>500</option><option>1000</option></select><button class="danger" onclick="clearLogs()">清空日志</button></div></div>
    <div class="tablewrap"><table><thead><tr><th>时间</th><th>触发</th><th>账号</th><th>Model</th><th>结果</th><th>HTTP</th><th>耗时</th><th>详情</th></tr></thead><tbody id="logRows"></tbody></table></div>
  </div>
</div>

<datalist id="modelSuggestions"><option value="gpt-5.6-luna"><option value="gpt-5.6-sol"><option value="gpt-5.6"><option value="gpt-6-astra"></datalist>

<dialog id="planDialog">
  <div class="modalhead" id="planDialogTitle">新建计划</div>
  <div class="modalbody">
    <input id="planId" type="hidden">
    <div class="grid2"><div class="field"><label>计划名称</label><input id="planName" placeholder="例如 06:00 点燃 A 账号"></div><div class="field"><label>账号</label><select id="planAccount"></select></div></div>
    <div class="grid2"><div class="field"><label>Model</label><input id="planModel" list="modelSuggestions"></div><div class="field"><label>时区</label><input id="planTimezone" placeholder="Asia/Shanghai"></div></div>
    <div class="field"><label>Cron 表达式</label><input id="planCron" class="mono" placeholder="0 6,11,16,21 * * *"><div class="hint">支持 *, */n, a-b, a-b/n 和逗号列表。</div></div>
    <div class="field"><label>Prompt</label><input id="planPrompt" placeholder="ping"></div>
    <label class="toggle"><input id="planEnabled" type="checkbox" checked> 启用此计划</label>
  </div>
  <div class="modalfoot"><button onclick="document.getElementById('planDialog').close()">取消</button><button class="primary" onclick="savePlan()">保存</button></div>
</dialog>

<dialog id="bulkPlanDialog">
  <div class="modalhead">一键为全部账号创建定时测试</div>
  <div class="modalbody">
    <div class="dangerbox">将为当前 CPA 中的每一个 Codex auth 账号创建一条独立定时计划，包括 disabled / unavailable 账号。已存在完全相同计划的账号会自动跳过，不会重复创建。</div>
    <div class="field"><label>计划名称前缀</label><input id="bulkPlanName" placeholder="全账号定时"></div>
    <div class="grid2"><div class="field"><label>Model</label><input id="bulkPlanModel" list="modelSuggestions"></div><div class="field"><label>时区</label><input id="bulkPlanTimezone" placeholder="Asia/Shanghai"></div></div>
    <div class="field"><label>Cron 表达式</label><input id="bulkPlanCron" class="mono" placeholder="0 6,11,16,21 * * *"><div class="hint">同一份时间表会应用到全部账号。例如每天 06:00 / 11:00 / 16:00 / 21:00：<span class="mono">0 6,11,16,21 * * *</span></div></div>
    <div class="field"><label>Prompt</label><input id="bulkPlanPrompt" placeholder="ping"></div>
    <label class="toggle"><input id="bulkPlanEnabled" type="checkbox" checked> 创建后立即启用全部计划</label>
    <div class="hint">创建后仍然是“一账号一计划”，因此你后续可以单独编辑、停用、立即发送或删除某个账号的计划。</div>
  </div>
  <div class="modalfoot"><button onclick="document.getElementById('bulkPlanDialog').close()">取消</button><button class="primary" onclick="createPlansForAll()">创建全部账号计划</button></div>
</dialog>

<div id="toast" class="toast"></div>
<script>
const BASE='/v0/management/plugins/cpa-scheduled-tests';
let STATE={accounts:[],plans:[],settings:{},jobs:[]};
let MODELS=['gpt-5.6-luna','gpt-5.6-sol','gpt-5.6','gpt-6-astra'];
function key(){return document.getElementById('mgmtKey').value.trim()}
function saveKey(){sessionStorage.setItem('cpa-scheduled-tests-key',key());toast('Management Key 已保存到当前标签页会话')}
function headers(json=true){const h={'Authorization':'Bearer '+key()};if(json)h['Content-Type']='application/json';return h}
async function api(path,opt={}){if(!key())throw new Error('请先输入 CPA Management Key');const r=await fetch(BASE+path,{...opt,headers:{...headers(opt.body!==undefined),...(opt.headers||{})}});const t=await r.text();let v;try{v=JSON.parse(t)}catch{v={error:t}}if(!r.ok)throw new Error(v.error||('HTTP '+r.status));return v}
async function hostApi(path,opt={}){if(!key())throw new Error('请先输入 CPA Management Key');const r=await fetch('/v0/management'+path,{...opt,headers:{...headers(opt.body!==undefined),...(opt.headers||{})}});const t=await r.text();let v;try{v=JSON.parse(t)}catch{v={error:t}}if(!r.ok)throw new Error(v.error||('HTTP '+r.status));return v}
async function loadModels(showToast=false){try{const r=await hostApi('/model-definitions/codex',{method:'GET'});const ids=(r.models||[]).map(x=>typeof x==='string'?x:x.id).filter(Boolean);if(ids.length)MODELS=[...new Set(ids)];document.getElementById('modelSuggestions').innerHTML=MODELS.map(x=>'<option value="'+esc(x)+'">').join('');if(showToast)toast('已同步 '+MODELS.length+' 个 Codex Model')}catch(e){if(showToast)toast('同步 Model 失败：'+e.message)}}
function toast(msg){const x=document.getElementById('toast');x.textContent=msg;x.style.display='block';clearTimeout(window.__tt);window.__tt=setTimeout(()=>x.style.display='none',3200)}
function esc(s){return String(s??'').replace(/[&<>'"]/g,c=>({'&':'&amp;','<':'&lt;','>':'&gt;',"'":'&#39;','"':'&quot;'}[c]))}
function fmtTime(s){if(!s||String(s).startsWith('0001-'))return '-';try{return new Date(s).toLocaleString()}catch{return s}}
function statusClass(s){return s==='success'?'ok':s==='limited'?'warn':s?'bad':'muted'}
function accountLabel(a){const n=a.label||a.email||a.name||a.id||a.auth_index;const flags=[a.disabled?'disabled':'',a.unavailable?'unavailable':''].filter(Boolean);return n+(flags.length?' ['+flags.join(', ')+']':'')}
function fillAccountSelect(id,value){const el=document.getElementById(id);el.innerHTML=STATE.accounts.map(a=>'<option value="'+esc(a.auth_index)+'">'+esc(accountLabel(a))+'</option>').join('');if(value)el.value=value}
async function loadAll(){try{const s=await api('/state',{method:'GET'});STATE=s;renderState();await loadModels(false);await loadLogs()}catch(e){toast(e.message)}}
function renderState(){document.getElementById('mAccounts').textContent=STATE.accounts.length;document.getElementById('mPlans').textContent=STATE.plans.length;document.getElementById('mEnabled').textContent=STATE.plans.filter(p=>p.enabled).length;document.getElementById('mRunning').textContent=STATE.jobs.filter(j=>j.running).length;fillAccountSelect('quickAccount');fillAccountSelect('planAccount');const set=STATE.settings||{};for(const id of ['quickModel','allModel'])if(!document.getElementById(id).value)document.getElementById(id).value=set.default_model||'gpt-5.6-luna';for(const id of ['quickPrompt','allPrompt'])if(!document.getElementById(id).value)document.getElementById(id).value=set.default_prompt||'ping';document.getElementById('bulkConcurrency').value=set.bulk_concurrency||4;renderPlans();renderJobs()}
function renderPlans(){const body=document.getElementById('planRows');if(!STATE.plans.length){body.innerHTML='<tr><td colspan="9" class="muted">暂无计划。点击“新建计划”。</td></tr>';return}body.innerHTML=STATE.plans.map(p=>'<tr><td>'+(p.enabled?'<span class="chip ok">启用</span>':'<span class="chip muted">停用</span>')+'</td><td>'+esc(p.name)+'</td><td class="ellipsis" title="'+esc(p.account_name||p.auth_index)+'">'+esc(p.account_name||p.auth_index)+'</td><td class="mono">'+esc(p.model)+'</td><td class="mono">'+esc(p.cron)+'</td><td>'+esc(p.timezone)+'</td><td class="nowrap">'+fmtTime(p.next_run)+'</td><td><span class="'+statusClass(p.last_status)+'">'+esc(p.last_status||'-')+'</span>'+(p.last_http_status?' · '+p.last_http_status:'')+'</td><td class="nowrap"><button class="small primary" onclick="runPlan(\''+esc(p.id)+'\')">立即发送</button> <button class="small" onclick="editPlan(\''+esc(p.id)+'\')">编辑</button> <button class="small danger" onclick="deletePlan(\''+esc(p.id)+'\')">删除</button></td></tr>').join('')}
function renderJobs(){const el=document.getElementById('jobs');if(!STATE.jobs.length){el.innerHTML='<span class="muted">暂无</span>';return}el.innerHTML=STATE.jobs.slice(0,10).map(j=>'<div class="job"><b>'+esc(j.kind)+'</b> · '+esc(j.model)+' · '+(j.running?'<span class="warn">运行中</span>':'<span class="muted">完成</span>')+' · '+j.completed+'/'+j.total+' · <span class="ok">'+j.succeeded+' 成功</span> / <span class="warn">'+j.limited+' 限流</span> / <span class="bad">'+j.failed+' 失败</span></div>').join('')}
function openPlan(p=null){document.getElementById('planDialogTitle').textContent=p?'编辑计划':'新建计划';document.getElementById('planId').value=p?.id||'';document.getElementById('planName').value=p?.name||'';fillAccountSelect('planAccount',p?.auth_index||STATE.accounts[0]?.auth_index||'');document.getElementById('planModel').value=p?.model||STATE.settings.default_model||'gpt-5.6-luna';document.getElementById('planTimezone').value=p?.timezone||STATE.settings.default_timezone||'Asia/Shanghai';document.getElementById('planCron').value=p?.cron||'0 6,11,16,21 * * *';document.getElementById('planPrompt').value=p?.prompt||STATE.settings.default_prompt||'ping';document.getElementById('planEnabled').checked=p?p.enabled:true;document.getElementById('planDialog').showModal()}
function editPlan(id){openPlan(STATE.plans.find(p=>p.id===id))}
function openBulkPlan(){const set=STATE.settings||{};document.getElementById('bulkPlanName').value='全账号定时';document.getElementById('bulkPlanModel').value=set.default_model||'gpt-5.6-luna';document.getElementById('bulkPlanTimezone').value=set.default_timezone||'Asia/Shanghai';document.getElementById('bulkPlanCron').value='0 6,11,16,21 * * *';document.getElementById('bulkPlanPrompt').value=set.default_prompt||'ping';document.getElementById('bulkPlanEnabled').checked=true;document.getElementById('bulkPlanDialog').showModal()}
async function createPlansForAll(){const body={name_prefix:document.getElementById('bulkPlanName').value,model:document.getElementById('bulkPlanModel').value,timezone:document.getElementById('bulkPlanTimezone').value,cron:document.getElementById('bulkPlanCron').value,prompt:document.getElementById('bulkPlanPrompt').value,enabled:document.getElementById('bulkPlanEnabled').checked};if(!body.model.trim()||!body.cron.trim()){toast('请填写 Model 和 Cron');return}if(!confirm('将为当前 '+STATE.accounts.length+' 个 Codex 账号批量创建定时计划。继续？'))return;try{const r=await api('/plans/bulk-create',{method:'POST',body:JSON.stringify(body)});document.getElementById('bulkPlanDialog').close();toast('批量创建完成：新增 '+r.created+'，跳过重复 '+r.skipped+'，账号总数 '+r.total_accounts);await loadAll()}catch(e){toast(e.message)}}
async function savePlan(){const p={id:document.getElementById('planId').value,name:document.getElementById('planName').value,auth_index:document.getElementById('planAccount').value,model:document.getElementById('planModel').value,timezone:document.getElementById('planTimezone').value,cron:document.getElementById('planCron').value,prompt:document.getElementById('planPrompt').value,enabled:document.getElementById('planEnabled').checked};try{await api('/plans/save',{method:'POST',body:JSON.stringify(p)});document.getElementById('planDialog').close();toast('计划已保存');await loadAll()}catch(e){toast(e.message)}}
async function deletePlan(id){if(!confirm('确定删除这个计划？'))return;try{await api('/plans/delete',{method:'POST',body:JSON.stringify({id})});toast('计划已删除');await loadAll()}catch(e){toast(e.message)}}
async function runPlan(id){try{const r=await api('/run',{method:'POST',body:JSON.stringify({plan_id:id})});toast('已启动立即发送：'+r.job_id);setTimeout(loadAll,800)}catch(e){toast(e.message)}}
async function quickRun(){try{const r=await api('/run',{method:'POST',body:JSON.stringify({auth_index:document.getElementById('quickAccount').value,model:document.getElementById('quickModel').value,prompt:document.getElementById('quickPrompt').value})});toast('已启动：'+r.job_id);setTimeout(loadAll,800)}catch(e){toast(e.message)}}
async function runAllAccounts(){const model=document.getElementById('allModel').value.trim();if(!model){toast('请填写 Model');return}if(!confirm('将向 CPA 中所有 Codex 账号强制发送一次真实请求，包括 disabled / unavailable。继续？'))return;try{const set={...STATE.settings,bulk_concurrency:Number(document.getElementById('bulkConcurrency').value||4)};await api('/settings',{method:'POST',body:JSON.stringify(set)});const r=await api('/run-all',{method:'POST',body:JSON.stringify({model,prompt:document.getElementById('allPrompt').value})});toast('全部账号任务已启动：'+r.job_id);setTimeout(loadAll,800)}catch(e){toast(e.message)}}
async function loadLogs(){if(!key())return;try{const n=document.getElementById('logLimit').value||200;const r=await api('/logs?limit='+encodeURIComponent(n),{method:'GET'});const body=document.getElementById('logRows');const logs=r.logs||[];if(!logs.length){body.innerHTML='<tr><td colspan="8" class="muted">暂无日志</td></tr>';return}body.innerHTML=logs.map(x=>'<tr><td class="nowrap">'+fmtTime(x.at)+'</td><td>'+esc(x.trigger)+'</td><td class="ellipsis" title="'+esc(x.account_name)+'">'+esc(x.account_name)+(x.disabled?' <span class="chip muted">disabled</span>':'')+(x.unavailable?' <span class="chip warn">unavailable</span>':'')+'</td><td class="mono">'+esc(x.model)+'</td><td class="'+statusClass(x.status)+'">'+esc(x.status)+'</td><td>'+(x.http_status||'-')+'</td><td>'+x.latency_ms+' ms</td><td class="ellipsis" title="'+esc(x.error||'')+'">'+esc(x.error||'-')+'</td></tr>').join('')}catch(e){toast(e.message)}}
async function clearLogs(){if(!confirm('确定清空插件运行日志？'))return;try{await api('/logs/clear',{method:'POST',body:'{}'});toast('日志已清空');await loadLogs()}catch(e){toast(e.message)}}
window.addEventListener('load',()=>{const k=sessionStorage.getItem('cpa-scheduled-tests-key')||'';document.getElementById('mgmtKey').value=k;if(k)loadAll()});
setInterval(()=>{if(key()&&STATE.jobs?.some(j=>j.running))loadAll()},2500);
</script>
</body></html>`
}
