#!/usr/bin/env python3
"""Verify a reviewed publication snapshot; not a legal clearance or signature check."""
from pathlib import Path, PurePosixPath
import hashlib
import json
import os
import re
import sys

ROOT = Path(__file__).resolve().parents[1]
MANIFEST = 'PUBLICATION-MANIFEST.json'
MAX_BYTES = 1024 * 1024
MAX_PHOTO_BYTES = 4 * 1024 * 1024
MAX_TOOL_BYTES = 12 * 1024 * 1024
REVIEWED_TOOLS = {'tools/flasher/GA36-Flasher-GUI-Windows-Linux.zip': '38e7b54f6f76aa5443ca43569550fda0d9e45bbcd3ae019b5b0bb673e60c6a0b', 'tools/flasher/screenshots/system.png': '5c28730b7e644efb21d9d7cbc7b5f00ccc27232dfa5f5f5ffb3b1f0750eb891e', 'tools/flasher/screenshots/apks.png': '9a1f33ed266e3e54e9ee095b90c1880d5684157d681bd12466d73c5b82a79609'}
REVIEWED_PHOTOS = {'assets/photos/android11-boot-logo.jpg': '61fc8723d461787bb703ee3545bf9de56638e2941765f878d88a373d7b871eac', 'assets/photos/android11-starting.jpg': 'd2ffe408bc329fbb43bf23e8d1796276aab9cfa45f374921cbf34c8787483920', 'assets/photos/android11-launcher.jpg': 'da714246dcfb7a8ea1d4e67ea9c9666e0aecfacda2d247e61dbb4dbcda9f7bea', 'assets/photos/android11-settings.jpg': 'c07fcd5d150a2f2099170657bbeec7eefe72d181dceb6a16bcd7e35cc4d6e3f0', 'assets/illustrations/system-selection.jpg': '56e7e115acd0f62133832cb2e4d680b9d2e09009e64d3bc90b9c48a2d4e69cd1', 'assets/photos/android7-device.jpg': '5ac482cb0387b04996f580b9925cc0843e92ba4bbe05d75a38fa4b8d8b1abcf4', 'assets/photos/android7-memory.jpg': 'd168cce2e68e85fd5e8b86aac026a472053507b9555908aa86627f7e9c87d10b', 'assets/photos/board-button-side.jpg': 'e0b17aa31e902d63bfe59c4b1ac964e5b8fe02aa20e3036088334fa1ea170663', 'assets/photos/board-component-side.jpg': '1104d16b5e6d4c18d74955036b5de96e40843becaf4fde68588e1133f67d5890'}
ALLOWED_SUFFIXES = {'.md', '.txt', '.c', '.py', '.kl', '.idc', '.json', '.go', '.mod', '.html', '.css', '.js', '.sha256'}
PATTERNS = [
    r'-----BEGIN [A-Z ]*PRIVATE KEY-----',
    r'github_pat_[A-Za-z0-9_]{20,}',
    r'gh[pousr]_[A-Za-z0-9]{20,}',
    r'AKIA[A-Z0-9]{16}',
    r'/home/[A-Za-z0-9._-]+/',
    r'/run/user/[0-9]+/',
    r'(?i)ro\.serialno\s*[=:]',
]

def inspect_text(path, rel, errors):
    if path.stat().st_size > MAX_BYTES:
        errors.append('Oversized file: ' + rel)
        return None
    data = path.read_bytes()
    if b'\0' in data:
        errors.append('Binary file: ' + rel)
        return None
    try:
        text = data.decode('utf-8')
    except UnicodeDecodeError:
        errors.append('Non-UTF8 file: ' + rel)
        return None
    for pattern in PATTERNS:
        if re.search(pattern, text):
            errors.append('Potential private value: ' + rel)
            break
    if path.suffix == '.json':
        try:
            json.loads(text)
        except (ValueError, RecursionError):
            errors.append('Invalid JSON: ' + rel)
    return data

def check():
    errors = []
    manifest_path = ROOT / MANIFEST
    if manifest_path.is_symlink() or not manifest_path.is_file():
        return ['Manifest missing or symlinked']
    raw = inspect_text(manifest_path, MANIFEST, errors)
    if raw is None or errors:
        return errors
    try:
        manifest = json.loads(raw)
        if manifest.get('schema_version') != 2:
            raise ValueError('Expected manifest schema version2')
        entries = manifest['files']
        if not isinstance(entries, list) or not entries:
            raise ValueError('Expected nonempty file list')
        expected = {}
        for entry in entries:
            rel = entry['path']
            path = PurePosixPath(rel)
            if (not isinstance(rel, str) or not rel or path.is_absolute()
                    or '..' in path.parts or str(path) != rel
                    or '\\' in rel or path.parts[0] == '.git'
                    or rel == MANIFEST or rel in expected):
                raise ValueError('Unsafe or duplicate manifest path')
            if not re.fullmatch(r'[0-9a-f]{64}', entry['sha256']):
                raise ValueError('Invalid SHA256')
            limit = MAX_TOOL_BYTES if rel in REVIEWED_TOOLS else (MAX_PHOTO_BYTES if rel in REVIEWED_PHOTOS else MAX_BYTES)
            if type(entry['bytes']) is not int or not 0 <= entry['bytes'] <= limit:
                raise ValueError('Invalid file size')
            expected[rel] = entry
    except (ValueError, TypeError, KeyError, AttributeError, RecursionError) as exc:
        return ['Invalid manifest: ' + str(exc)]

    actual = set()
    def walk_error(exc):
        raise exc
    for current, dirs, files in os.walk(ROOT, followlinks=False, onerror=walk_error):
        directory = Path(current)
        # Git metadata is local and is never part of a publication snapshot.
        if directory == ROOT:
            if '.git' in dirs and not (ROOT / '.git').is_symlink():
                dirs.remove('.git')
            files = [name for name in files if name != '.git' or (ROOT / name).is_symlink()]
        for name in list(dirs):
            child = directory / name
            if child.is_symlink():
                errors.append('Symlink is forbidden: ' + child.relative_to(ROOT).as_posix())
                dirs.remove(name)
        for name in files:
            p = directory / name
            rel = p.relative_to(ROOT).as_posix()
            if p.is_symlink():
                errors.append('Symlink is forbidden: ' + rel)
                continue
            if not p.is_file():
                errors.append('Nonregular file: ' + rel)
                continue
            actual.add(rel)
            if rel == MANIFEST:
                continue
            if rel not in expected:
                errors.append('Unreviewed file: ' + rel)
                continue
            if p.suffix.lower() not in ALLOWED_SUFFIXES and rel != '.gitignore' and rel not in REVIEWED_PHOTOS and rel not in REVIEWED_TOOLS:
                errors.append('Unexpected file type: ' + rel)
            if rel in REVIEWED_TOOLS:
                if p.stat().st_size > MAX_TOOL_BYTES:
                    errors.append('Oversized reviewed tool: ' + rel)
                    continue
                data = p.read_bytes()
                if hashlib.sha256(data).hexdigest() != REVIEWED_TOOLS[rel]:
                    errors.append('Tool differs from reviewed artifact: ' + rel)
                if rel.endswith('.png') and not data.startswith(b'\x89PNG\r\n\x1a\n'):
                    errors.append('Invalid PNG framing: ' + rel)
                if rel.endswith('.zip'):
                    import zipfile
                    try:
                        with zipfile.ZipFile(p) as z:
                            if z.testzip() is not None:
                                errors.append('Corrupt GUI package: ' + rel)
                            for info in z.infolist():
                                name = PurePosixPath(info.filename)
                                if name.is_absolute() or '..' in name.parts or info.filename.lower().endswith(('.img', '.img.gz', '.iso', '.ext4', '.apk', '.obb', '.cpio', '.so', '.ko')):
                                    errors.append('Forbidden package content: ' + info.filename)
                    except (OSError, zipfile.BadZipFile):
                        errors.append('Invalid GUI ZIP: ' + rel)
            elif rel in REVIEWED_PHOTOS:
                if p.stat().st_size > MAX_PHOTO_BYTES:
                    errors.append('Oversized reviewed photo: ' + rel)
                    continue
                data = p.read_bytes()
                if not data.startswith(b'\xff\xd8') or not data.endswith(b'\xff\xd9'):
                    errors.append('Invalid JPEG framing: ' + rel)
                if hashlib.sha256(data).hexdigest() != REVIEWED_PHOTOS[rel]:
                    errors.append('Photo differs from explicitly reviewed image: ' + rel)
            else:
                data = inspect_text(p, rel, errors)
            if data is None:
                continue
            if len(data) != expected[rel]['bytes']:
                errors.append('Changed size: ' + rel)
            if hashlib.sha256(data).hexdigest() != expected[rel]['sha256']:
                errors.append('Changed file: ' + rel)
    for rel in sorted(set(expected) - actual):
        errors.append('Missing file: ' + rel)
    return errors

if __name__ == '__main__':
    try:
        problems = check()
    except (OSError, ValueError) as exc:
        problems = ['Check failed: ' + str(exc)]
    if problems:
        print('\n'.join(problems))
        sys.exit(1)
    print('PASS: reviewed file set, UTF8/JSON text and explicitly allowlisted JPEG/PNG images and GUI-only tool package, sizes and SHA256 match; no unreviewed files or symlinks.')
    print('Local root .git metadata excluded. Gitignore rules do not authorize extra files.')
    print('The manifest is not signed; this check is not a legal or exhaustive secret-scan guarantee.')
