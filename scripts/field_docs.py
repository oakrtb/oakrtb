"""Generate field reference and schema descriptions; reject stale translations."""
import argparse
import json
import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
PROTO = ROOT / 'proto/oakrtb/v2/openrtb.proto'
SCHEMA = ROOT / 'schema/jsonschema/openrtb.schema.json'
NATIVE = ROOT / 'schema/jsonschema/native.schema.json'
TEXT = ROOT / 'docs/field-descriptions.json'
OUTPUT = ROOT / 'docs/fields.md'


def parse_proto(text, kind='message'):
    blocks = {}
    for match in re.finditer(r'^' + kind + r'\s+(\w+)\s*\{\n(.*?)^}', text, re.M | re.S):
        fields, comments = {}, []
        for line in match[2].splitlines():
            stripped = line.strip()
            if stripped.startswith('//'):
                comments.append(stripped[2:].strip())
                continue
            field = re.match(r'\s*(?:(optional|repeated)\s+)?(\w+)\s+(\w+)\s*=\s*(\d+)', line)
            enum = re.match(r'\s*(\w+)\s*=\s*(-?\d+)', line) if kind == 'enum' else None
            if field or enum:
                if '//' in line:
                    comments.append(line.split('//', 1)[1].strip())
                if field:
                    fields[field[3]] = {'type': ' '.join(filter(None, field.group(1, 2))), 'description': ' '.join(comments)}
                else:
                    fields[enum[1]] = {'type': enum[2], 'description': ' '.join(comments)}
            comments = []
        blocks[match[1]] = fields
    return blocks


def native_objects(schema):
    objects = {}

    def visit(name, node):
        if 'properties' in node:
            objects[name] = node
            for key, prop in node['properties'].items():
                visit(name + '.' + key, prop)
        if isinstance(node.get('items'), dict):
            visit(name + '[]', node['items'])

    visit('NativeRequest', schema)
    for name, node in schema.get('$defs', {}).items():
        visit('Native' + name, node)
    return objects


def validate_descriptions(fields, descriptions):
    expected = {f'{obj}.{name}': data for obj, values in fields.items() for name, data in values.items()}
    if expected.keys() != descriptions.keys():
        raise ValueError(f'Translation coverage mismatch: missing={expected.keys() - descriptions.keys()}, extra={descriptions.keys() - expected.keys()}')
    for key, data in expected.items():
        if not data['description'] or not descriptions[key].get('zh', '').strip():
            raise ValueError(f'Missing description: {key}')
        if descriptions[key]['source'] != data['description']:
            raise ValueError(f'Stale translation: {key}; review zh and update source in docs/field-descriptions.json')


def cell(value):
    return str(value).replace('|', '&#124;').replace('\n', ' ')


def constraints(node):
    return json.dumps({k: v for k, v in node.items() if k not in ('description', 'properties', 'title', '$defs')}, ensure_ascii=False, separators=(',', ':'))


def render(proto, schema, native, translations):
    fields = parse_proto(proto)
    validate_descriptions(fields, translations['fields'])
    lines = ['# 完整字段手册', '',
             '此文件由 `python3 scripts/field_docs.py --write` 生成，请勿直接编辑。中文释义维护于 `docs/field-descriptions.json`。', '',
             '覆盖当前 proto 的全部消息字段和枚举，以及 Native.request 内嵌 JSON 的全部 Schema 属性。', '',
             '## 阅读约定', '',
             '- JSON 列中的“是”只表示该对象的 Schema `required`；对象自身是否必须出现由父对象决定。跨字段约束另列。',
             '- protobuf 的 `optional` 表示保留字段存在性，不等于业务可选；`repeated` 对应数组。protobuf 本身不执行 JSON Schema 校验。',
             '- JSON 约束直接取自 Schema；`$ref` 指向文末基础类型或本手册同名对象。释义中的标准建议不代表 SDK 已强制验证。',
             '- `ext` 在 JSON 中是对象，在 protobuf 中保存为 JSON 字符串；Native.request 在两种外层编码中均为 JSON 字符串。',
             '- 枚举表用于解释值；int32 字段的实际允许范围以 JSON Schema 为准，不能仅凭枚举表推断拒绝未知值。', '',
             '## 业务规则与尺寸', '',
             '- Banner：Readiness 要求正数 w/h 成对提供，或 format 非空；这不是 Schema 的 required 规则。format 的每个元素仍应表达有效固定尺寸或比例尺寸，完整合同校验请调用 jsonschema。',
             '- Bid.mtype：对应 Imp 同时包含多种素材形态时，BidCheck 要求显式提供；Schema 只在字段出现时限制为 1–4。',
             '- Banner.w/h、Format.w/h/wmin、Video.w/h、Bid.w/h 使用设备无关像素（DIPs）；Device.w/h 使用屏幕物理像素。pxratio = 物理像素 / DIPs。广告容器尺寸不等于整块屏幕尺寸。',
             '- Native 图片宽高是图片资源的像素尺寸，不应统一换成广告容器的 DIPs。Imp.rwdd 表示激励广告，不限视频。',
             '- 开屏可结合 Imp.instl 描述全屏/插屏；Banner 是素材形态；信息流可结合 Native.request.plcmttype 描述。这些不是同一个字段中的互斥版位枚举。', '',
             '尺寸语义参照 [IAB OpenRTB 2.6](https://github.com/InteractiveAdvertisingBureau/openrtb2.x/blob/main/2.6.md)。校验层职责见 [SDK 架构](sdk.md)，协议规则见 [协议规范](spec.md)。', '',
             '## 对象索引', '', ' · '.join(f'[{name}](#{name.lower()})' for name in fields), '']
    for obj, values in fields.items():
        definition = schema['$defs'][obj]
        if values.keys() != definition['properties'].keys():
            raise ValueError(f'Proto/schema field mismatch: {obj}')
        lines += [f'## {obj}', '', '| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |', '|---|---|---|---|---|']
        for name, data in values.items():
            prop = definition['properties'][name]
            prop['description'] = data['description']
            lines.append(f"| `{name}` | `{data['type']}` | {'是' if name in definition.get('required', []) else '否'} | `{cell(constraints(prop))}` | {cell(translations['fields'][obj + '.' + name]['zh'])} |")
        rules = {k: v for k, v in definition.items() if k in ('allOf', 'anyOf', 'oneOf', 'not', 'dependentRequired', 'if', 'then', 'else')}
        if rules:
            lines += ['', '对象级 JSON 约束：', '', '```json', json.dumps(rules, ensure_ascii=False, indent=2), '```']
        lines.append('')
    lines += ['## Schema 基础类型', '', '```json', json.dumps({k: v for k, v in schema['$defs'].items() if k not in fields}, ensure_ascii=False, indent=2), '```', '',
              '## Native.request 内嵌 JSON', '',
              '以下为本仓库 Native 请求 Schema 的支持范围，不包含 Native 响应完整协议。先解析外层 request 字符串，再按 native.schema.json 校验；外层字符串类型检查不等于内层合同校验。', '',
              '枚举语义参考 [IAB Native 1.2](https://github.com/InteractiveAdvertisingBureau/Native-Ads/blob/main/OpenRTB-Native-Ads-Specification-Final-1.2.md)。当前 Schema 对多数分类字段只验证整数类型，未实现完整标准的枚举和组合约束。', '']
    native_fields = {f'{obj}.{name}' for obj, node in native_objects(native).items() for name in node['properties']}
    if native_fields != translations['native'].keys():
        raise ValueError(f'Native translation coverage mismatch: {native_fields ^ translations["native"].keys()}')
    for obj, node in native_objects(native).items():
        lines += [f'### {obj}', '', '| 字段 | JSON 必填 | JSON 约束 | 中文说明 |', '|---|---|---|---|']
        for name, prop in node['properties'].items():
            desc = translations['native'][f'{obj}.{name}']
            if not desc.strip():
                raise ValueError(f'Missing Native description: {obj}.{name}')
            lines.append(f"| `{name}` | {'是' if name in node.get('required', []) else '否'} | `{cell(constraints(prop))}` | {cell(desc)} |")
        lines.append('')
    lines += ['## protobuf 枚举参考', '', '以下保留源码的英文枚举说明；0 值可能只是 protobuf 未设置哨兵，是否可用于 JSON 取决于字段约束。', '']
    for name, values in parse_proto(proto, 'enum').items():
        lines += [f'### {name}', '', '| 名称 | 值 | 说明 |', '|---|---|---|']
        lines += [f"| `{key}` | {value['type']} | {cell(value['description'])} |" for key, value in values.items()]
        lines.append('')
    return '\n'.join(lines)


def schema_text(original, schema):
    # Preserve layout and numeric literals; replace or append only property descriptions.
    pattern = re.compile(r'"(?:[^"\\]|\\.)*"|[{}]', re.S)
    stack, edits, pending = [], [], None
    for token in pattern.finditer(original):
        value = token.group()
        if value.startswith('"'):
            if original[token.end():].lstrip().startswith(':'):
                pending = json.loads(value)
        elif value == '{':
            stack.append((pending, token.end()))
            pending = None
        else:
            path = [key for key, _ in stack]
            if len(path) == 5 and path[1] == '$defs' and path[3] == 'properties':
                obj, field = path[2], path[4]
                desc = schema['$defs'].get(obj, {}).get('properties', {}).get(field, {}).get('description')
                if desc is not None:
                    start = stack[-1][1]
                    body = original[start:token.start()]
                    old = None
                    depth = 0
                    for part in pattern.finditer(body):
                        if part.group() == '{':
                            depth += 1
                        elif part.group() == '}':
                            depth -= 1
                        elif depth == 0 and part.group() == '"description"':
                            old = re.compile(r'"description"\s*:\s*"(?:[^"\\]|\\.)*"').match(body, part.start())
                            if old:
                                break
                    entry = '"description": ' + json.dumps(desc, ensure_ascii=False)
                    if old:
                        edits.append((start + old.start(), start + old.end(), entry))
                    else:
                        edits.append((start, start, ' ' + entry + ','))
            stack.pop()
            pending = None
    for start, end, replacement in reversed(edits):
        original = original[:start] + replacement + original[end:]
    if json.loads(original) != schema:
        raise ValueError('Schema description update failed to preserve the structure')
    return original


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    mode = parser.add_mutually_exclusive_group(required=True)
    mode.add_argument('--write', action='store_true')
    mode.add_argument('--check', action='store_true')
    args = parser.parse_args()
    original = SCHEMA.read_text()
    schema = json.loads(original)
    rendered = render(PROTO.read_text(), schema, json.loads(NATIVE.read_text()), json.loads(TEXT.read_text()))
    outputs = {OUTPUT: rendered, SCHEMA: schema_text(original, schema)}
    for path, expected in outputs.items():
        if args.write:
            path.write_text(expected)
        elif not path.exists() or path.read_text() != expected:
            raise SystemExit(f'Stale generated file: {path.relative_to(ROOT)}; run python3 scripts/field_docs.py --write')
    print('Field documentation and schema descriptions are synchronized.')


if __name__ == '__main__':
    main()
