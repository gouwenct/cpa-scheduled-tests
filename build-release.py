"""Build a native CPA library and ZIP, or verify all five ZIPs and write checksums."""
import hashlib
import os
from pathlib import Path
import re
import subprocess
import sys
import zipfile

ROOT = Path(__file__).resolve().parent
os.chdir(ROOT)
VERSION = re.search(r'pluginVersion\s*=\s*"([0-9]+\.[0-9]+\.[0-9]+)"',
                    (ROOT / 'main.go').read_text(encoding='utf-8')).group(1)
PLUGIN = 'cpa-scheduled-tests'
PLATFORMS = {
    ('darwin', 'amd64'): 'dylib', ('darwin', 'arm64'): 'dylib',
    ('linux', 'amd64'): 'so', ('linux', 'arm64'): 'so',
    ('windows', 'amd64'): 'dll',
}
DIST = ROOT / 'dist'
DIST.mkdir(exist_ok=True)

if os.environ.get('GITHUB_REF', '').startswith('refs/tags/'):
    if os.environ['GITHUB_REF'] != f'refs/tags/v{VERSION}':
        raise SystemExit('Release tag must match pluginVersion in main.go')

if sys.argv[1:] == ['--checksums']:
    archives = []
    for (goos, goarch), extension in PLATFORMS.items():
        archive = DIST / f'{PLUGIN}_{VERSION}_{goos}_{goarch}.zip'
        with zipfile.ZipFile(archive) as bundle:
            if bundle.namelist() != [f'{PLUGIN}.{extension}']:
                raise SystemExit(f'Invalid library layout: {archive.name}')
            if bundle.testzip() is not None:
                raise SystemExit(f'Corrupt ZIP: {archive.name}')
        archives.append(archive)
elif not sys.argv[1:]:
    goos, goarch = subprocess.check_output(['go', 'env', 'GOOS', 'GOARCH'], text=True).split()
    extension = PLATFORMS.get((goos, goarch))
    if extension is None:
        raise SystemExit(f'Unsupported platform: {goos}/{goarch}')
    library = DIST / f'{PLUGIN}.{extension}'
    env = dict(os.environ, CGO_ENABLED='1')
    subprocess.run(['go', 'build', '-trimpath', '-buildmode=c-shared',
                    '-ldflags=-s -w', '-o', str(library), '.'], env=env, check=True)
    # Load the native library with a fake CPA host; no real credentials or requests.
    subprocess.run([sys.executable, 'abi_smoke_test.py', str(library)], check=True)
    archive = DIST / f'{PLUGIN}_{VERSION}_{goos}_{goarch}.zip'
    with zipfile.ZipFile(archive, 'w', zipfile.ZIP_DEFLATED) as bundle:
        bundle.write(library, arcname=library.name)
    archives = [archive]
else:
    raise SystemExit('Usage: python build-release.py [--checksums]')

lines = [f'{hashlib.sha256(p.read_bytes()).hexdigest()}  {p.name}\n' for p in archives]
(DIST / 'checksums.txt').write_text(''.join(lines), encoding='utf-8')
print(''.join(lines), end='')
