#!/usr/bin/env python3
"""Check the reviewed source/docs snapshot. This is not a legal clearance tool."""
from pathlib import Path
import hashlib,json,re,sys
root=Path(__file__).resolve().parents[1]
manifest=json.loads((root/'PUBLICATION-MANIFEST.json').read_text())
expected={item['path']:item for item in manifest['files']}
actual=set()
errors=[]
for p in root.rglob('*'):
 rel=p.relative_to(root).as_posix()
 if p.is_symlink():errors.append('Symlink is forbidden: '+rel);continue
 if p.is_dir():continue
 actual.add(rel)
 if rel=='PUBLICATION-MANIFEST.json':continue
 if rel not in expected:errors.append('Unreviewed file: '+rel);continue
 data=p.read_bytes()
 if len(data)>1024*1024 or b'\0' in data:errors.append('Binary/oversized file: '+rel);continue
 try:text=data.decode('utf-8')
 except UnicodeDecodeError:errors.append('Non-UTF8 file: '+rel);continue
 if hashlib.sha256(data).hexdigest()!=expected[rel]['sha256']:errors.append('Changed file: '+rel)
 if len(data)!=expected[rel]['bytes']:errors.append('Changed size: '+rel)
 if p.suffix.lower() not in {'.md','.txt','.c','.py','.kl','.idc'} and p.name!='.gitignore':errors.append('Unexpected file type: '+rel)
 if rel!='tools/check-publication.py':
  patterns=[r'-----BEGIN [A-Z ]*PRIVATE KEY-----',r'github_pat_[A-Za-z0-9_]{20,}',r'gh[pousr]_[A-Za-z0-9]{20,}',r'AKIA[A-Z0-9]{16}',r'/home/test/',r'/run/user/',r'R3CTA0X201P',r'(?i)ro\.serialno\s*[=:]']
  for pattern in patterns:
   if re.search(pattern,text):errors.append('Potential private value: '+rel)
for rel in sorted(set(expected)-actual):errors.append('Missing file: '+rel)
if errors:
 print('\n'.join(errors));sys.exit(1)
print(f'PASS: {len(expected)} reviewed text/source files; hashes match; no extra files or symlinks.')
print('No firmware images, Android binaries, games, card images or private keys in this reviewed snapshot.')
print('Re-run after changes. Git ignore rules alone do not prevent sensitive uploads.')
