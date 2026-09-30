# 完整字段手册

此文件由 `python3 scripts/field_docs.py --write` 生成，请勿直接编辑。中文释义维护于 `docs/field-descriptions.json`。

覆盖当前 proto 的全部消息字段和枚举，以及 Native.request 内嵌 JSON 的全部 Schema 属性。

## 阅读约定

- JSON 列中的“是”只表示该对象的 Schema `required`；对象自身是否必须出现由父对象决定。跨字段约束另列。
- protobuf 的 `optional` 表示保留字段存在性，不等于业务可选；`repeated` 对应数组。protobuf 本身不执行 JSON Schema 校验。
- JSON 约束直接取自 Schema；`$ref` 指向文末基础类型或本手册同名对象。释义中的标准建议不代表 SDK 已强制验证。
- `ext` 在 JSON 中是对象，在 protobuf 中保存为 JSON 字符串；Native.request 在两种外层编码中均为 JSON 字符串。
- 枚举表用于解释值；int32 字段的实际允许范围以 JSON Schema 为准，不能仅凭枚举表推断拒绝未知值。

## 业务规则与尺寸

- Banner：Readiness 要求正数 w/h 成对提供，或 format 非空；这不是 Schema 的 required 规则。format 的每个元素仍应表达有效固定尺寸或比例尺寸，完整合同校验请调用 jsonschema。
- Bid.mtype：对应 Imp 同时包含多种素材形态时，BidCheck 要求显式提供；Schema 只在字段出现时限制为 1–4。
- Banner.w/h、Format.w/h/wmin、Video.w/h、Bid.w/h 使用设备无关像素（DIPs）；Device.w/h 使用屏幕物理像素。pxratio = 物理像素 / DIPs。广告容器尺寸不等于整块屏幕尺寸。
- Native 图片宽高是图片资源的像素尺寸，不应统一换成广告容器的 DIPs。Imp.rwdd 表示激励广告，不限视频。
- 开屏可结合 Imp.instl 描述全屏/插屏；Banner 是素材形态；信息流可结合 Native.request.plcmttype 描述。这些不是同一个字段中的互斥版位枚举。

尺寸语义参照 [IAB OpenRTB 2.6](https://github.com/InteractiveAdvertisingBureau/openrtb2.x/blob/main/2.6.md)。校验层职责见 [SDK 架构](sdk.md)，协议规则见 [协议规范](spec.md)。

## 对象索引

[BidRequest](#bidrequest) · [Imp](#imp) · [Metric](#metric) · [Banner](#banner) · [Format](#format) · [Video](#video) · [Audio](#audio) · [Native](#native) · [Pmp](#pmp) · [Deal](#deal) · [Qty](#qty) · [DurFloors](#durfloors) · [Refresh](#refresh) · [RefSettings](#refsettings) · [Site](#site) · [App](#app) · [Dooh](#dooh) · [Publisher](#publisher) · [Producer](#producer) · [Content](#content) · [Network](#network) · [Channel](#channel) · [Device](#device) · [UserAgent](#useragent) · [BrandVersion](#brandversion) · [Geo](#geo) · [User](#user) · [Data](#data) · [Segment](#segment) · [EID](#eid) · [UID](#uid) · [Source](#source) · [SupplyChain](#supplychain) · [SupplyChainNode](#supplychainnode) · [Regs](#regs) · [BidResponse](#bidresponse) · [SeatBid](#seatbid) · [Bid](#bid)

## BidRequest

| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|---|
| `id` | `string` | 是 | `{"$ref":"#/$defs/NonEmptyString"}` | 本场拍卖唯一 ID。必填。BidResponse.id 必须原样回传，用于对账与日志关联。 |
| `imp` | `repeated Imp` | 是 | `{"type":"array","minItems":1,"items":{"$ref":"#/$defs/Imp"}}` | 待售展示机会列表。必填，至少 1 个。一单可多 Imp（多广告位/多形态）。 |
| `site` | `Site` | 否 | `{"$ref":"#/$defs/Site"}` | 网站库存。与 app、dooh 互斥。场景：网页/移动 Web 广告位。 |
| `app` | `App` | 否 | `{"$ref":"#/$defs/App"}` | 应用库存。与 site、dooh 互斥。场景：iOS/Android App 内广告。 |
| `dooh` | `Dooh` | 否 | `{"$ref":"#/$defs/Dooh"}` | 数字户外库存。与 site、app 互斥。场景：DOOH 屏/看板。 |
| `device` | `Device` | 否 | `{"$ref":"#/$defs/Device"}` | 设备与环境。强烈推荐。场景：定向、反作弊、创意适配（尺寸/OS）。 |
| `user` | `User` | 否 | `{"$ref":"#/$defs/User"}` | 用户/受众。推荐。场景：频控、人群定向、同意信号（consent/eids）。 |
| `test` | `optional int32` | 否 | `{"type":"integer","enum":[0,1]}` | 测试模式。取值 FlagBool：0=生产，1=测试。场景：联调；测试流量不得计入计费。 |
| `at` | `optional int32` | 是 | `{"allOf":[{"$ref":"#/$defs/AuctionType"},{"anyOf":[{"enum":[1,2]},{"minimum":500}],"description":"Required: 1 first price, 2 second price plus, or >=500 exchange-specific. Fixed price (3) is only valid for Deal.at."}]}` | 拍卖类型，必填：1=一价，2=二价加价，≥500=交易平台自定义；3=固定价格仅适用于 Deal.at。Schema 与构建器拒绝未指定的 0。 |
| `tmax` | `optional int32` | 否 | `{"type":"integer","minimum":1}` | 超时毫秒（含网络）。场景：DSP 必须在 tmax 内返回，否则 Exchange 当 204。 |
| `wseat` | `repeated string` | 否 | `{"type":"array","items":{"type":"string"}}` | 白名单 seat。场景：只允许名单内买家席位参竞。 |
| `bseat` | `repeated string` | 否 | `{"type":"array","items":{"type":"string"}}` | 黑名单 seat。场景：排除特定买家席位。 |
| `allimps` | `optional int32` | 否 | `{"type":"integer","enum":[0,1]}` | 是否带齐全部 Imp。取值 FlagBool：0=否，1=是。场景：帮助 DSP 做全局预算与互斥决策。 |
| `cur` | `repeated string` | 是 | `{"type":"array","minItems":1,"items":{"$ref":"#/$defs/Iso4217"}}` | 可接受出价币种（ISO-4217）。必填，至少 1 个。场景：多币种结算；未列出则 DSP 勿用该币种。 |
| `wlang` | `repeated string` | 否 | `{"type":"array","items":{"type":"string"}}` | 允许的创意语言（ISO-639-1）。场景：语言定向过滤。 |
| `wlangb` | `repeated string` | 否 | `{"type":"array","items":{"type":"string"}}` | 允许的创意语言（BCP-47）。场景：比 wlang 更细的语言标签。 |
| `acat` | `repeated string` | 否 | `{"type":"array","items":{"type":"string"}}` | 允许的广告主行业类目。场景：白名单类目库存。 |
| `bcat` | `repeated string` | 否 | `{"type":"array","items":{"type":"string"}}` | 屏蔽的广告主行业类目。场景：品牌安全/类目黑名单。 |
| `cattax` | `optional int32` | 否 | `{"type":"integer"}` | 类目 taxonomy。取值 CategoryTaxonomy：1–7 见枚举；>=500 厂商自定义。场景：解释 cat/bcat/acat。 |
| `badv` | `repeated string` | 否 | `{"type":"array","items":{"type":"string"}}` | 屏蔽广告主域名。场景：竞品/敏感域名屏蔽。 |
| `bapp` | `repeated string` | 否 | `{"type":"array","items":{"type":"string"}}` | 屏蔽应用 bundle/包名。场景：App 流量上屏蔽特定广告主 App。 |
| `source` | `Source` | 否 | `{"$ref":"#/$defs/Source"}` | 上游来源与供应链。场景：schain 透明化、header bidding 链路审计。 |
| `regs` | `Regs` | 否 | `{"$ref":"#/$defs/Regs"}` | 法规与隐私信号。场景：COPPA/GDPR/CCPA/GPP 合规决策。 |
| `ext` | `string` | 否 | `{"$ref":"#/$defs/Ext"}` | 扩展：JSON 对象字符串。场景：Exchange 私有字段，双方约定后使用。 |

对象级 JSON 约束：

```json
{
  "not": {
    "anyOf": [
      {
        "required": [
          "site",
          "app"
        ]
      },
      {
        "required": [
          "site",
          "dooh"
        ]
      },
      {
        "required": [
          "app",
          "dooh"
        ]
      }
    ]
  }
}
```

## Imp

| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|---|
| `id` | `string` | 是 | `{"$ref":"#/$defs/NonEmptyString"}` | Imp 在本请求内唯一。必填。Bid.impid 必须回填此值。 |
| `metric` | `repeated Metric` | 否 | `{"type":"array","items":{"$ref":"#/$defs/Metric"}}` | 库存质量/可验证指标。场景：viewability、视频完成率等第三方度量。 |
| `banner` | `Banner` | 否 | `{"$ref":"#/$defs/Banner"}` | 展示广告规格。场景：横幅、富媒体、插屏图片等。 |
| `video` | `Video` | 否 | `{"$ref":"#/$defs/Video"}` | 视频广告规格。场景：前贴/中贴/后贴、激励视频、CTV pod。 |
| `audio` | `Audio` | 否 | `{"$ref":"#/$defs/Audio"}` | 音频广告规格。场景：播客、电台流、音乐 App。 |
| `native` | `Native` | 否 | `{"$ref":"#/$defs/Native"}` | 原生广告规格。场景：信息流卡片；request 内嵌 Native 1.2 JSON。 |
| `pmp` | `Pmp` | 否 | `{"$ref":"#/$defs/Pmp"}` | 私有交易。场景：PDB/PD/Preferred Deal，带 deals[]。 |
| `displaymanager` | `string` | 否 | `{"type":"string"}` | 广告渲染 SDK/播放器名。场景：识别 Mediation/Player，做兼容定向。 |
| `displaymanagerver` | `string` | 否 | `{"type":"string"}` | displaymanager 版本。场景：按 SDK 版本灰度或排障。 |
| `instl` | `optional int32` | 否 | `{"type":"integer","enum":[0,1]}` | 是否插屏/全屏。取值 FlagBool：0=否，1=是。场景：插屏比价与体验策略。 |
| `tagid` | `string` | 否 | `{"type":"string"}` | 发布商广告位 ID。场景：按位优化、报表、底价策略。 |
| `bidfloor` | `optional double` | 否 | `{"type":"number","minimum":0}` | CPM 底价。场景：低于此价的出价通常无效。 |
| `bidfloorcur` | `string` | 否 | `{"$ref":"#/$defs/Iso4217"}` | 底价币种（ISO-4217）。场景：与 bidfloor 成对；默认常 USD。 |
| `clickbrowser` | `optional int32` | 否 | `{"type":"integer","enum":[0,1]}` | 点击打开方式。取值 FlagBool：0=嵌入 WebView，1=独立浏览器。场景：移动端落地页与归因。 |
| `secure` | `optional int32` | 否 | `{"type":"integer","enum":[0,1]}` | 是否要求 HTTPS。取值 FlagBool：0=否，1=创意与页面须 HTTPS。场景：混合内容拦截环境。 |
| `iframebuster` | `repeated string` | 否 | `{"type":"array","items":{"type":"string"}}` | 支持的 iframe buster 厂商名。场景：可跳出 iframe 的富媒体。 |
| `rwdd` | `optional int32` | 否 | `{"type":"integer","enum":[0,1]}` | 激励广告标记（FlagBool）：0=无奖励，1=观看广告后获得奖励；适用于多种素材形态，不限于视频。 |
| `ssai` | `optional int32` | 否 | `{"type":"integer","enum":[0,1,2,3]}` | 服务端广告插入。取值 Ssai：0=未知，1=客户端，2=服务端，3=混合。场景：CTV/直播 SSAI。 |
| `exp` | `optional int32` | 否 | `{"type":"integer","minimum":0}` | 展示过期秒数建议。场景：长时间会话中创意缓存多久失效。 |
| `qty` | `Qty` | 否 | `{"$ref":"#/$defs/Qty"}` | 数量乘数。场景：DOOH/成组曝光，用 multiplier 换算计费展示量。 |
| `dt` | `optional double` | 否 | `{"type":"number"}` | 距展示的估计秒数。场景：提前竞价（early auction）时的时间窗。 |
| `refresh` | `Refresh` | 否 | `{"$ref":"#/$defs/Refresh"}` | 自动刷新详情。场景：页内广告位自动轮换。 |
| `ext` | `string` | 否 | `{"$ref":"#/$defs/Ext"}` | 扩展 JSON 字符串。 |

对象级 JSON 约束：

```json
{
  "anyOf": [
    {
      "required": [
        "banner"
      ]
    },
    {
      "required": [
        "video"
      ]
    },
    {
      "required": [
        "audio"
      ]
    },
    {
      "required": [
        "native"
      ]
    }
  ]
}
```

## Metric

| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|---|
| `type` | `string` | 是 | `{"$ref":"#/$defs/NonEmptyString"}` | 指标类型，如 "viewability"、"completion_rate"。 |
| `value` | `optional double` | 是 | `{"type":"number"}` | 指标取值，通常 0–1 概率或比率。 |
| `vendor` | `string` | 否 | `{"type":"string"}` | 度量提供方。场景：IAS、MOAT、发布商自测等。 |
| `ext` | `string` | 否 | `{"$ref":"#/$defs/Ext"}` | 扩展：JSON 对象字符串。 |

## Banner

| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|---|
| `format` | `repeated Format` | 否 | `{"type":"array","items":{"$ref":"#/$defs/Format"}}` | 可接受尺寸集合。场景：弹性/多尺寸位，优先于单独的 w/h。 |
| `w` | `optional int32` | 否 | `{"type":"integer","minimum":1}` | 固定展示宽度，单位为设备无关像素（DIPs），与 h 配对。 |
| `h` | `optional int32` | 否 | `{"type":"integer","minimum":1}` | 固定展示高度，单位为设备无关像素（DIPs），与 w 配对。 |
| `btype` | `repeated int32` | 否 | `{"type":"array","items":{"type":"integer"}}` | 屏蔽的横幅类型。取值 BannerAdType：1=XHTML文本，2=XHTML Banner，3=JS，4=iframe。 |
| `battr` | `repeated int32` | 否 | `{"type":"array","items":{"type":"integer"}}` | 屏蔽的创意属性。取值 CreativeAttribute：1–18 见枚举；>=500 厂商自定义。场景：禁自动播放声等。 |
| `pos` | `optional int32` | 否 | `{"type":"integer"}` | 广告位置。取值 AdPosition：0=未知，1=首屏上，2=锁定，3=首屏下，4=头，5=脚，6=侧栏，7=全屏。 |
| `mimes` | `repeated string` | 否 | `{"type":"array","items":{"type":"string"}}` | 支持的 MIME。场景：image/jpeg、image/png、text/javascript 等。 |
| `topframe` | `optional int32` | 否 | `{"type":"integer","enum":[0,1]}` | 是否顶层 frame。取值 FlagBool：0=在 iframe 内，1=顶层。场景：是否允许跳出/expand。 |
| `expdir` | `repeated int32` | 否 | `{"type":"array","items":{"type":"integer"}}` | 允许扩展方向。取值 ExpandableDirection：1=左，2=右，3=上，4=下，5=全屏，6=缩小。 |
| `api` | `repeated int32` | 否 | `{"type":"array","items":{"type":"integer"}}` | 支持的 API。取值 ApiFramework：1=VPAID1，2=VPAID2，3=MRAID1，4=ORMMA，5=MRAID2，6=MRAID3，7=OMID，8–9=SIMID；>=500 厂商。 |
| `id` | `string` | 否 | `{"type":"string"}` | 同伴广告时区分 Banner 的 ID。场景：视频 companion 多 banner 时关联。 |
| `vcm` | `optional int32` | 否 | `{"type":"integer","enum":[0,1]}` | 视频伴生模式。取值 VideoCompanionMode：0=与视频同时，1=结束后 end card。 |
| `ext` | `string` | 否 | `{"$ref":"#/$defs/Ext"}` | 扩展：JSON 对象字符串。 |

## Format

| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|---|
| `w` | `optional int32` | 否 | `{"type":"integer","minimum":1}` | 允许的展示宽度，单位为设备无关像素（DIPs）。 |
| `h` | `optional int32` | 否 | `{"type":"integer","minimum":1}` | 允许的展示高度，单位为设备无关像素（DIPs）。 |
| `wratio` | `optional int32` | 否 | `{"type":"integer","minimum":1}` | 宽比例；与 hratio 描述纵横比 |
| `hratio` | `optional int32` | 否 | `{"type":"integer","minimum":1}` | 高比例 |
| `wmin` | `optional int32` | 否 | `{"type":"integer","minimum":1}` | 按宽高比例布局时的最小展示宽度，单位为设备无关像素（DIPs）。 |
| `ext` | `string` | 否 | `{"$ref":"#/$defs/Ext"}` | 扩展：JSON 对象字符串。 |

## Video

| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|---|
| `mimes` | `repeated string` | 是 | `{"type":"array","minItems":1,"items":{"type":"string"}}` | 支持的视频 MIME。必填。场景：video/mp4、video/webm 等。 |
| `minduration` | `optional int32` | 否 | `{"type":"integer","minimum":0}` | 最短时长（秒）。场景：过滤过短素材。 |
| `maxduration` | `optional int32` | 否 | `{"type":"integer","minimum":0}` | 最长时长（秒）。场景：过滤过长素材。 |
| `startdelay` | `optional int32` | 否 | `{"type":"integer"}` | 起播延迟（秒）。特殊值：0=前贴，-1=通用中贴，-2=后贴，>0=中贴且延迟该秒数。场景：贴片位置。 |
| `maxseq` | `optional int32` | 否 | `{"type":"integer","minimum":1}` | pod 内最大广告条数。场景：CTV ad pod 容量。 |
| `poddur` | `optional int32` | 否 | `{"type":"integer","minimum":1}` | 整个 pod 总时长（秒）。场景：多条广告共享的时间预算。 |
| `protocols` | `repeated int32` | 否 | `{"type":"array","items":{"type":"integer"}}` | 支持的视频协议。取值 Protocol：1–3=VAST1–3，4–6=对应 Wrapper，7–8=VAST4，9–10=DAAST，11–16=VAST4.1–4.3及Wrapper。 |
| `w` | `optional int32` | 否 | `{"type":"integer","minimum":1}` | 播放器宽度，单位为设备无关像素（DIPs）。 |
| `h` | `optional int32` | 否 | `{"type":"integer","minimum":1}` | 播放器高度，单位为设备无关像素（DIPs）。 |
| `podid` | `string` | 否 | `{"type":"string"}` | 广告 pod ID。场景：同一动态 pod 内多 Imp 共享，用于编排。 |
| `podseq` | `optional int32` | 否 | `{"type":"integer"}` | 该 Imp 所属 pod 序号策略。取值 PodSequence：0=任意，1=第一个；特殊 -1=最后一个 pod。 |
| `rqddurs` | `repeated int32` | 否 | `{"type":"array","items":{"type":"integer"}}` | 要求的精确时长列表（秒）。场景：只接受 15s/30s 等标准时长。 |
| `plcmt` | `optional int32` | 否 | `{"type":"integer"}` | 视频放置类型（替代旧 placement）。取值 VideoPlcmt：1=Instream，2=伴随内容，3=插屏，4=无内容，5=暂停，6=屏保，7=叠层，8=分屏，9=画面植入。 |
| `linearity` | `optional int32` | 否 | `{"type":"integer"}` | 线性模式。取值 VideoLinearity：1=线性贴片，2=非线性/Overlay。 |
| `skip` | `optional int32` | 否 | `{"type":"integer","enum":[0,1]}` | 是否可跳过。取值 FlagBool：0=不可跳，1=可跳。场景：可跳过前贴 vs 不可跳。 |
| `skipmin` | `optional int32` | 否 | `{"type":"integer","minimum":0}` | 允许跳过前的最少视频秒数。 |
| `skipafter` | `optional int32` | 否 | `{"type":"integer","minimum":0}` | 出现跳过按钮前的秒数。 |
| `slotinpod` | `optional int32` | 否 | `{"type":"integer"}` | pod 内槽位。取值 SlotInPod：0=任意，1=第一，2=第一或最后，3=第一/中/最后；特殊 -1=最后槽。场景：包段选槽。 |
| `mincpmpersec` | `optional double` | 否 | `{"type":"number"}` | 每秒最低 CPM。场景：按时长计价的底价约束。 |
| `battr` | `repeated int32` | 否 | `{"type":"array","items":{"type":"integer"}}` | 屏蔽的创意属性。取值 CreativeAttribute：1–18 见枚举；>=500 厂商自定义。 |
| `maxextended` | `optional int32` | 否 | `{"type":"integer"}` | 允许超出 maxduration 的最大秒数；-1 表示不限。 |
| `minbitrate` | `optional int32` | 否 | `{"type":"integer","minimum":0}` | 最低码率（Kbps）。 |
| `maxbitrate` | `optional int32` | 否 | `{"type":"integer","minimum":0}` | 最高码率（Kbps）。 |
| `boxingallowed` | `optional int32` | 否 | `{"type":"integer","enum":[0,1]}` | 是否允许 letterboxing。取值 FlagBool：0=否，1=是。场景：画幅不匹配时黑边。 |
| `playbackmethod` | `repeated int32` | 否 | `{"type":"array","items":{"type":"integer"}}` | 播放方式。取值 PlaybackMethod：1=加载有声，2=加载静音，3=点击有声，4=悬停有声，5=入视口有声，6=入视口静音，7=连续播。 |
| `playbackend` | `optional int32` | 否 | `{"type":"integer"}` | 播放终止方式。取值 PlaybackCessationMode：1=播完/用户停，2=离视口/用户停，3=离视口后悬浮至播完。 |
| `delivery` | `repeated int32` | 否 | `{"type":"array","items":{"type":"integer"}}` | 传输方式。取值 DeliveryMethod：1=流式，2=渐进，3=下载。 |
| `pos` | `optional int32` | 否 | `{"type":"integer"}` | 广告位置。取值 AdPosition（同 Banner.pos）。 |
| `companionad` | `repeated Banner` | 否 | `{"type":"array","items":{"$ref":"#/$defs/Banner"}}` | 伴生 Banner 规格。场景：视频旁/下方 companion。 |
| `api` | `repeated int32` | 否 | `{"type":"array","items":{"type":"integer"}}` | 支持的 API。取值 ApiFramework（见枚举）。 |
| `companiontype` | `repeated int32` | 否 | `{"type":"array","items":{"type":"integer"}}` | 伴生广告类型。取值 CompanionType：1=Static，2=HTML，3=iframe。 |
| `placement` | `optional int32` | 否 | `{"type":"integer"}` | 已弃用的 placement；请用 plcmt（OpenRTB 2.6-202303+）。 |
| `poddedupe` | `repeated int32` | 否 | `{"type":"array","items":{"type":"integer"}}` | Pod 去重策略（OpenRTB 2.6-202402+）。取值见 AdCOM Pod Deduplication。 Pod 去重策略（OpenRTB 2.6-202402+）。取值 PodDeduplication：1=按 adomain，2=按 IAB Content Taxonomy，3=按创意 ID，4=按 mediafile URL；可多选。 |
| `durfloors` | `repeated DurFloors` | 否 | `{"type":"array","items":{"$ref":"#/$defs/DurFloors"}}` | 按时长分段底价。 按时长分段底价。场景：15s/30s 等不同时长档使用不同 CPM 底价。 |
| `ext` | `string` | 否 | `{"$ref":"#/$defs/Ext"}` | 扩展：JSON 对象字符串。 |

## Audio

| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|---|
| `mimes` | `repeated string` | 是 | `{"type":"array","minItems":1,"items":{"type":"string"}}` | 支持的音频 MIME。必填。 |
| `minduration` | `optional int32` | 否 | `{"type":"integer","minimum":0}` | 最短秒数 |
| `maxduration` | `optional int32` | 否 | `{"type":"integer","minimum":0}` | 最长秒数 |
| `poddur` | `optional int32` | 否 | `{"type":"integer","minimum":1}` | pod 总时长 |
| `protocols` | `repeated int32` | 否 | `{"type":"array","items":{"type":"integer"}}` | 音频协议。取值 Protocol（含 DAAST 9/10、VAST 等，见枚举）。 |
| `startdelay` | `optional int32` | 否 | `{"type":"integer"}` | 起播延迟。同 Video：0=前贴，-1=中贴，-2=后贴，>0=中贴延迟秒数。 |
| `rqddurs` | `repeated int32` | 否 | `{"type":"array","items":{"type":"integer"}}` | 要求的精确时长 |
| `podid` | `string` | 否 | `{"type":"string"}` | 音频 pod ID |
| `podseq` | `optional int32` | 否 | `{"type":"integer"}` | pod 序号。取值 PodSequence：0=任意，1=第一；特殊 -1=最后 pod。 |
| `slotinpod` | `optional int32` | 否 | `{"type":"integer"}` | pod 槽位。取值 SlotInPod；特殊 -1=最后槽。 |
| `mincpmpersec` | `optional double` | 否 | `{"type":"number"}` | 每秒最低 CPM。场景：按时长计价的底价约束。 |
| `battr` | `repeated int32` | 否 | `{"type":"array","items":{"type":"integer"}}` | 屏蔽创意属性。取值 CreativeAttribute。 |
| `maxextended` | `optional int32` | 否 | `{"type":"integer"}` | 允许超出 maxduration 的最大秒数；-1 表示不限。 |
| `minbitrate` | `optional int32` | 否 | `{"type":"integer"}` | 最低码率（Kbps）。 |
| `maxbitrate` | `optional int32` | 否 | `{"type":"integer"}` | 最高码率（Kbps）。 |
| `delivery` | `repeated int32` | 否 | `{"type":"array","items":{"type":"integer"}}` | 传输方式。取值 DeliveryMethod：1=流式，2=渐进，3=下载。 |
| `companionad` | `repeated Banner` | 否 | `{"type":"array","items":{"$ref":"#/$defs/Banner"}}` | 可视化伴生 |
| `api` | `repeated int32` | 否 | `{"type":"array","items":{"type":"integer"}}` | 支持的 API。取值 ApiFramework。 |
| `companiontype` | `repeated int32` | 否 | `{"type":"array","items":{"type":"integer"}}` | 伴生类型。取值 CompanionType：1=Static，2=HTML，3=iframe。 |
| `maxseq` | `optional int32` | 否 | `{"type":"integer"}` | pod 内最大广告数 |
| `feed` | `optional int32` | 否 | `{"type":"integer"}` | 音频流类型。取值 FeedType：1=音乐服务，2=FM/AM，3=播客。 |
| `stitched` | `optional int32` | 否 | `{"type":"integer","enum":[0,1]}` | 是否已与内容缝合。取值 FlagBool：0=否，1=是。 |
| `nvol` | `optional int32` | 否 | `{"type":"integer"}` | 音量归一。取值 VolumeNormalizationMode：0=无，1=平均，2=峰值，3=响度，4=自定义。 |
| `durfloors` | `repeated DurFloors` | 否 | `{"type":"array","items":{"$ref":"#/$defs/DurFloors"}}` | 按时长分段底价。 |
| `ext` | `string` | 否 | `{"$ref":"#/$defs/Ext"}` | 扩展：JSON 对象字符串。 |

## Native

| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|---|
| `request` | `string` | 是 | `{"type":"string","minLength":1}` | Native 请求体（JSON 字符串）。必填。内含 assets[] 等；需按 Native 1.2 校验。 |
| `ver` | `string` | 否 | `{"type":"string"}` | Native 规范版本，如 "1.2"。场景：告诉 DSP 解析哪一版 markup。 |
| `api` | `repeated int32` | 否 | `{"type":"array","items":{"type":"integer"}}` | 支持的 API。取值 ApiFramework（见枚举）。 |
| `battr` | `repeated int32` | 否 | `{"type":"array","items":{"type":"integer"}}` | 屏蔽的创意属性。取值 CreativeAttribute（见枚举）。 |
| `ext` | `string` | 否 | `{"$ref":"#/$defs/Ext"}` | 扩展：JSON 对象字符串。 |

## Pmp

| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|---|
| `private_auction` | `optional int32` | 否 | `{"type":"integer","enum":[0,1]}` | 是否仅私有拍卖。取值 FlagBool：0=可对公开市场，1=仅私有。场景：PDB 封闭交易。 |
| `deals` | `repeated Deal` | 否 | `{"type":"array","items":{"$ref":"#/$defs/Deal"}}` | 可用成交条件列表。场景：一条 Imp 挂多个 Deal 供 DSP 选择。 |
| `ext` | `string` | 否 | `{"$ref":"#/$defs/Ext"}` | 扩展：JSON 对象字符串。 |

## Deal

| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|---|
| `id` | `string` | 是 | `{"$ref":"#/$defs/NonEmptyString"}` | Deal ID。必填。Bid.dealid 赢价时回填。 |
| `bidfloor` | `optional double` | 否 | `{"type":"number","minimum":0}` | 该 Deal 底价（CPM）。 |
| `bidfloorcur` | `string` | 否 | `{"$ref":"#/$defs/Iso4217"}` | 底价币种。 |
| `at` | `optional int32` | 否 | `{"$ref":"#/$defs/AuctionType"}` | 该 Deal 拍卖类型。取值 AuctionType：1=一价，2=二价+，3=固定价；>=500 自定义。场景：覆盖 BidRequest.at。 |
| `wseat` | `repeated string` | 否 | `{"type":"array","items":{"type":"string"}}` | 允许的买家 seat 白名单。 |
| `wadomain` | `repeated string` | 否 | `{"type":"array","items":{"type":"string"}}` | 允许的广告主域名白名单。 |
| `guar` | `optional int32` | 否 | `{"type":"integer","enum":[0,1]}` | 是否保量。取值 FlagBool：0=否，1=保量 Deal。场景：程序化合约保量。 |
| `mincpmpersec` | `optional double` | 否 | `{"type":"number"}` | 每秒最低 CPM（视频按时长）。 |
| `durfloors` | `repeated DurFloors` | 否 | `{"type":"array","items":{"$ref":"#/$defs/DurFloors"}}` | 按时长分段底价（视频/音频 deal）。 |
| `ext` | `string` | 否 | `{"$ref":"#/$defs/Ext"}` | 扩展：JSON 对象字符串。 |

## Qty

| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|---|
| `multiplier` | `optional double` | 否 | `{"type":"number"}` | 数量乘数（相对 1 次标准展示）。 |
| `sourcetype` | `optional int32` | 否 | `{"type":"integer"}` | 乘数来源。取值 QtySourceType：0=度量厂商方法，1=发布商/厂商特定方法。 |
| `vendor` | `string` | 否 | `{"type":"string"}` | 提供乘数的厂商。 |
| `ext` | `string` | 否 | `{"$ref":"#/$defs/Ext"}` | 扩展：JSON 对象字符串。 |

## DurFloors

| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|---|
| `mindur` | `optional int32` | 否 | `{"type":"integer","minimum":0}` | 区间最短秒数（含） |
| `maxdur` | `optional int32` | 否 | `{"type":"integer","minimum":0}` | 区间最长秒数（含）；可省略表示无上界 |
| `bidfloor` | `optional double` | 否 | `{"type":"number"}` | 该时长区间 CPM 底价 |
| `bidfloorcur` | `string` | 否 | `{"$ref":"#/$defs/Iso4217"}` | 底价币种 ISO-4217 |
| `ext` | `string` | 否 | `{"$ref":"#/$defs/Ext"}` | 扩展：JSON 对象字符串。 |

## Refresh

| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|---|
| `refsettings` | `repeated RefSettings` | 否 | `{"type":"array","items":{"$ref":"#/$defs/RefSettings"}}` | 刷新触发与间隔 |
| `count` | `optional int32` | 否 | `{"type":"integer"}` | 自上次整页加载以来已刷新次数 |
| `ext` | `string` | 否 | `{"$ref":"#/$defs/Ext"}` | 扩展：JSON 对象字符串。 |

## RefSettings

| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|---|
| `reftype` | `optional int32` | 否 | `{"type":"integer"}` | 自动刷新触发类型。取值 AutoRefreshTrigger：0=未知，1=用户动作，2=事件驱动，3=按时间间隔。 |
| `minint` | `optional int32` | 否 | `{"type":"integer"}` | 最小刷新间隔（秒）。场景：限制刷新频率。 |
| `ext` | `string` | 否 | `{"$ref":"#/$defs/Ext"}` | 扩展：JSON 对象字符串。 |

## Site

| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|---|
| `id` | `string` | 否 | `{"type":"string"}` | Exchange 侧站点 ID |
| `name` | `string` | 否 | `{"type":"string"}` | 站点名（可展示名） |
| `domain` | `string` | 否 | `{"type":"string"}` | 主域，如 example.com；品牌安全与屏蔽 |
| `cattax` | `optional int32` | 否 | `{"type":"integer"}` | 类目 taxonomy。取值 CategoryTaxonomy（见枚举） |
| `cat` | `repeated string` | 否 | `{"type":"array","items":{"type":"string"}}` | 站点类目 |
| `sectioncat` | `repeated string` | 否 | `{"type":"array","items":{"type":"string"}}` | 栏目类目 |
| `pagecat` | `repeated string` | 否 | `{"type":"array","items":{"type":"string"}}` | 当前页类目 |
| `page` | `string` | 否 | `{"type":"string"}` | 完整页面 URL；上下文定向 |
| `ref` | `string` | 否 | `{"type":"string"}` | Referrer URL |
| `search` | `string` | 否 | `{"type":"string"}` | 搜索词（若来自搜索结果页） |
| `mobile` | `optional int32` | 否 | `{"type":"integer","enum":[0,1]}` | 是否移动优化站。取值 FlagBool：0=否，1=是 |
| `privacypolicy` | `optional int32` | 否 | `{"type":"integer","enum":[0,1]}` | 是否有隐私政策。取值 FlagBool：0=否，1=是 |
| `publisher` | `Publisher` | 否 | `{"$ref":"#/$defs/Publisher"}` | 发布商 |
| `content` | `Content` | 否 | `{"$ref":"#/$defs/Content"}` | 页内内容元数据 |
| `keywords` | `string` | 否 | `{"type":"string"}` | 逗号分隔关键词（兼容旧字段） |
| `kwarray` | `repeated string` | 否 | `{"type":"array","items":{"type":"string"}}` | 关键词数组（推荐） |
| `inventorypartnerdomain` | `string` | 否 | `{"type":"string"}` | ads.txt/库存合作域 |
| `ext` | `string` | 否 | `{"$ref":"#/$defs/Ext"}` | 扩展：JSON 对象字符串。 |

## App

| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|---|
| `id` | `string` | 否 | `{"type":"string"}` | Exchange 侧 App ID |
| `name` | `string` | 否 | `{"type":"string"}` | 应用名 |
| `bundle` | `string` | 否 | `{"type":"string"}` | 包名/Bundle ID；应用定向与屏蔽核心字段 |
| `domain` | `string` | 否 | `{"type":"string"}` | 应用关联域 |
| `storeurl` | `string` | 否 | `{"type":"string"}` | 应用商店详情 URL |
| `cattax` | `optional int32` | 否 | `{"type":"integer"}` | 类目 taxonomy。取值 CategoryTaxonomy |
| `cat` | `repeated string` | 否 | `{"type":"array","items":{"type":"string"}}` | 应用类目。场景：行业定向；taxonomy 见 cattax。 |
| `sectioncat` | `repeated string` | 否 | `{"type":"array","items":{"type":"string"}}` | 栏目/分区类目。 |
| `pagecat` | `repeated string` | 否 | `{"type":"array","items":{"type":"string"}}` | 当前视图/页面类目。 |
| `ver` | `string` | 否 | `{"type":"string"}` | 应用版本 |
| `privacypolicy` | `optional int32` | 否 | `{"type":"integer","enum":[0,1]}` | 是否有隐私政策。取值 FlagBool |
| `paid` | `optional int32` | 否 | `{"type":"integer","enum":[0,1]}` | 是否付费应用。取值 FlagBool：0=免费，1=付费 |
| `publisher` | `Publisher` | 否 | `{"$ref":"#/$defs/Publisher"}` | 发布商。场景：应用归属的发行商主体。 |
| `content` | `Content` | 否 | `{"$ref":"#/$defs/Content"}` | 应用内内容元数据。场景：内嵌播放器/文章上下文。 |
| `keywords` | `string` | 否 | `{"type":"string"}` | 逗号分隔关键词（兼容旧字段）。 |
| `kwarray` | `repeated string` | 否 | `{"type":"array","items":{"type":"string"}}` | 关键词数组（推荐）。 |
| `inventorypartnerdomain` | `string` | 否 | `{"type":"string"}` | 库存合作域（app-ads.txt 合作方）。 |
| `ext` | `string` | 否 | `{"$ref":"#/$defs/Ext"}` | 扩展：JSON 对象字符串。 |

## Dooh

| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|---|
| `id` | `string` | 否 | `{"type":"string"}` | 屏点/单元 ID |
| `name` | `string` | 否 | `{"type":"string"}` | 屏点/单元名称。 |
| `venue` | `optional int32` | 否 | `{"type":"integer"}` | 旧版数值场所类型；优先用 venuetype + venuetypetax |
| `fixed` | `optional int32` | 否 | `{"type":"integer"}` | 是否固定屏。取值 FlagBool：0=可移动，1=固定 |
| `publisher` | `Publisher` | 否 | `{"$ref":"#/$defs/Publisher"}` | 发布商/屏网运营商。 |
| `domain` | `string` | 否 | `{"type":"string"}` | 关联业务域。 |
| `keywords` | `string` | 否 | `{"type":"string"}` | 逗号分隔关键词（兼容旧字段）。 |
| `kwarray` | `repeated string` | 否 | `{"type":"array","items":{"type":"string"}}` | 关键词数组（推荐）。 |
| `content` | `Content` | 否 | `{"$ref":"#/$defs/Content"}` | 当前播放的内容上下文 |
| `venuetype` | `repeated string` | 否 | `{"type":"array","items":{"type":"string"}}` | 场所类型 ID 列表（OpenOOH 等）；taxonomy 由 venuetypetax 指定。 |
| `venuetypetax` | `optional int32` | 否 | `{"type":"integer"}` | 场所 taxonomy。AdCOM DOOH Venue Taxonomies；默认 1。 |
| `ext` | `string` | 否 | `{"$ref":"#/$defs/Ext"}` | 扩展：JSON 对象字符串。 |

## Publisher

| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|---|
| `id` | `string` | 否 | `{"type":"string"}` | 发布商 ID；报表与合约主体 |
| `name` | `string` | 否 | `{"type":"string"}` | 发布商名称。 |
| `cattax` | `optional int32` | 否 | `{"type":"integer"}` | 类目 taxonomy。取值 CategoryTaxonomy |
| `cat` | `repeated string` | 否 | `{"type":"array","items":{"type":"string"}}` | 发布商类目 |
| `domain` | `string` | 否 | `{"type":"string"}` | 发布商主域 |
| `ext` | `string` | 否 | `{"$ref":"#/$defs/Ext"}` | 扩展：JSON 对象字符串。 |

## Producer

| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|---|
| `id` | `string` | 否 | `{"type":"string"}` | 内容制片方/出品方 |
| `name` | `string` | 否 | `{"type":"string"}` | 出品方名称。 |
| `cattax` | `optional int32` | 否 | `{"type":"integer"}` | 类目 taxonomy。取值 CategoryTaxonomy |
| `cat` | `repeated string` | 否 | `{"type":"array","items":{"type":"string"}}` | 出品方类目；taxonomy 见 cattax。 |
| `domain` | `string` | 否 | `{"type":"string"}` | 出品方主域。 |
| `ext` | `string` | 否 | `{"$ref":"#/$defs/Ext"}` | 扩展：JSON 对象字符串。 |

## Content

| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|---|
| `id` | `string` | 否 | `{"type":"string"}` | 内容 ID。场景：剧集/文章唯一标识，用于上下文定向。 |
| `episode` | `optional int32` | 否 | `{"type":"integer"}` | 集数 |
| `title` | `string` | 否 | `{"type":"string"}` | 标题 |
| `series` | `string` | 否 | `{"type":"string"}` | 系列名 |
| `season` | `string` | 否 | `{"type":"string"}` | 季 |
| `artist` | `string` | 否 | `{"type":"string"}` | 艺术家（音频/音乐） |
| `genre` | `string` | 否 | `{"type":"string"}` | 自由文本类型（旧字段）；有 taxonomy 时优先 genres+gtax |
| `album` | `string` | 否 | `{"type":"string"}` | 专辑名（音频/音乐）。 |
| `isrc` | `string` | 否 | `{"type":"string"}` | 录音制品国际标准码 |
| `producer` | `Producer` | 否 | `{"$ref":"#/$defs/Producer"}` | 内容出品方。 |
| `url` | `string` | 否 | `{"type":"string"}` | 内容 URL |
| `cattax` | `optional int32` | 否 | `{"type":"integer"}` | 类目 taxonomy。取值 CategoryTaxonomy |
| `cat` | `repeated string` | 否 | `{"type":"array","items":{"type":"string"}}` | 内容类目；taxonomy 见 cattax。 |
| `prodq` | `optional int32` | 否 | `{"type":"integer"}` | 制作质量。取值 ProductionQuality：0=未知，1=专业，2=准专业，3=UGC |
| `context` | `optional int32` | 否 | `{"type":"integer"}` | 内容上下文。取值 ContentContext：1=视频，2=游戏，3=音乐，4=应用，5=文本，6=其他，7=未知 |
| `contentrating` | `string` | 否 | `{"type":"string"}` | 内容分级，如 MPAA |
| `userrating` | `string` | 否 | `{"type":"string"}` | 用户评分 |
| `qagmediarating` | `optional int32` | 否 | `{"type":"integer"}` | 媒体评级。取值 MediaRating：1=全年龄，2=12+，3=成人 |
| `keywords` | `string` | 否 | `{"type":"string"}` | 逗号分隔关键词（兼容旧字段）。 |
| `kwarray` | `repeated string` | 否 | `{"type":"array","items":{"type":"string"}}` | 关键词数组（推荐）。 |
| `livestream` | `optional int32` | 否 | `{"type":"integer","enum":[0,1]}` | 播出方式（OpenRTB 2.6-202606）：0=非排期（VOD/用户发起），1=排期/线性；≠ realtime。 |
| `sourcerelationship` | `optional int32` | 否 | `{"type":"integer"}` | 与发布商关系。取值 SourceRelationship：0=间接，1=直接 |
| `len` | `optional int32` | 否 | `{"type":"integer"}` | 内容时长秒数 |
| `language` | `string` | 否 | `{"type":"string"}` | ISO-639-1 |
| `langb` | `string` | 否 | `{"type":"string"}` | BCP-47 |
| `embeddable` | `optional int32` | 否 | `{"type":"integer","enum":[0,1]}` | 是否可嵌入。取值 FlagBool：0=否，1=是 |
| `data` | `repeated Data` | 否 | `{"type":"array","items":{"$ref":"#/$defs/Data"}}` | 内容侧分段数据 |
| `network` | `Network` | 否 | `{"$ref":"#/$defs/Network"}` | 播出网络 |
| `channel` | `Channel` | 否 | `{"$ref":"#/$defs/Channel"}` | 频道 |
| `gtax` | `optional int32` | 否 | `{"type":"integer"}` | 类型 taxonomy（OpenRTB 2.6-202501+）；缺省按 IAB 指南可视为 Content Taxonomy 3.1。 类型 taxonomy（OpenRTB 2.6-202501+）。解释 genres[]；缺省可按 IAB Content Taxonomy 约定。 |
| `genres` | `repeated string` | 否 | `{"type":"array","items":{"type":"string"}}` | 类型 ID 列表，taxonomy 由 gtax 定义（OpenRTB 2.6-202501+）。 |
| `realtime` | `optional int32` | 否 | `{"type":"integer","enum":[0,1]}` | 是否实时发生（OpenRTB 2.6-202606）：0=非实时（如回放），1=实时（如直播赛事）。 |
| `firstbroadcast` | `optional int32` | 否 | `{"type":"integer","enum":[0,1]}` | 是否首播（OpenRTB 2.6-202606）：0=非首次，1=首次对观众播出。 |
| `ext` | `string` | 否 | `{"$ref":"#/$defs/Ext"}` | 扩展：JSON 对象字符串。 |

## Network

| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|---|
| `id` | `string` | 否 | `{"type":"string"}` | 网络 ID。 |
| `name` | `string` | 否 | `{"type":"string"}` | 网络名称。 |
| `domain` | `string` | 否 | `{"type":"string"}` | 网络主域。 |
| `ext` | `string` | 否 | `{"$ref":"#/$defs/Ext"}` | 扩展：JSON 对象字符串。 |

## Channel

| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|---|
| `id` | `string` | 否 | `{"type":"string"}` | 频道 ID。 |
| `name` | `string` | 否 | `{"type":"string"}` | 频道名称。 |
| `domain` | `string` | 否 | `{"type":"string"}` | 频道主域。 |
| `ext` | `string` | 否 | `{"$ref":"#/$defs/Ext"}` | 扩展：JSON 对象字符串。 |

## Device

| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|---|
| `geo` | `Geo` | 否 | `{"$ref":"#/$defs/Geo"}` | 设备地理位置（可与 User.geo 不同） |
| `dnt` | `optional int32` | 否 | `{"type":"integer","enum":[0,1]}` | Do Not Track。取值 FlagBool：0=否，1=开启 DNT |
| `lmt` | `optional int32` | 否 | `{"type":"integer","enum":[0,1]}` | 限制广告跟踪（LAT）。取值 FlagBool：0=否，1=限制 |
| `ua` | `string` | 否 | `{"type":"string"}` | User-Agent 原始串；解析失败时的兜底 |
| `sua` | `UserAgent` | 否 | `{"$ref":"#/$defs/UserAgent"}` | 结构化 UA（SUA）；优先于 ua 做精准定向 |
| `ip` | `string` | 否 | `{"type":"string"}` | IPv4；地理与风控 |
| `ipv6` | `string` | 否 | `{"type":"string"}` | IPv6 |
| `devicetype` | `optional int32` | 否 | `{"type":"integer"}` | 设备类型。取值 DeviceType：1=通用移动，2=PC，3=CTV，4=手机，5=平板，6=联网设备，7=机顶盒，8=OOH |
| `make` | `string` | 否 | `{"type":"string"}` | 厂商，如 Apple |
| `model` | `string` | 否 | `{"type":"string"}` | 型号 |
| `os` | `string` | 否 | `{"type":"string"}` | 操作系统 |
| `osv` | `string` | 否 | `{"type":"string"}` | OS 版本 |
| `hwv` | `string` | 否 | `{"type":"string"}` | 硬件版本 |
| `h` | `optional int32` | 否 | `{"type":"integer","minimum":0}` | 屏幕物理高度，单位为物理像素；不是广告容器高度。 |
| `w` | `optional int32` | 否 | `{"type":"integer","minimum":0}` | 屏幕物理宽度，单位为物理像素；不是广告容器宽度。 |
| `ppi` | `optional int32` | 否 | `{"type":"integer"}` | 像素密度 |
| `pxratio` | `optional double` | 否 | `{"type":"number"}` | 物理像素与设备无关像素（DIPs）的比值。 |
| `js` | `optional int32` | 否 | `{"type":"integer","enum":[0,1]}` | 是否支持 JS。取值 FlagBool：0=否，1=是 |
| `geofetch` | `optional int32` | 否 | `{"type":"integer","enum":[0,1]}` | 是否允许创意再取地理。取值 FlagBool：0=否，1=是 |
| `language` | `string` | 否 | `{"type":"string"}` | 设备语言 ISO-639-1 |
| `langb` | `string` | 否 | `{"type":"string"}` | BCP-47 |
| `carrier` | `string` | 否 | `{"type":"string"}` | 运营商 |
| `mccmnc` | `string` | 否 | `{"type":"string"}` | MCC-MNC |
| `connectiontype` | `optional int32` | 否 | `{"type":"integer"}` | 连接类型。取值 ConnectionType：0=未知，1=以太网，2=WiFi，3=蜂窝未知，4=2G，5=3G，6=4G，7=5G |
| `ifa` | `string` | 否 | `{"type":"string"}` | 广告 ID（IDFA/GAID 等）；用户级定向与归因 |
| `ext` | `string` | 否 | `{"$ref":"#/$defs/Ext"}` | 扩展：JSON 对象字符串。 |

## UserAgent

| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|---|
| `browsers` | `repeated BrandVersion` | 否 | `{"type":"array","items":{"$ref":"#/$defs/BrandVersion"}}` | 浏览器品牌版本列表 |
| `platform` | `BrandVersion` | 否 | `{"$ref":"#/$defs/BrandVersion"}` | 平台（OS）品牌版本 |
| `mobile` | `optional int32` | 否 | `{"type":"integer","enum":[0,1]}` | 是否移动。取值 FlagBool：0=否，1=是 |
| `architecture` | `string` | 否 | `{"type":"string"}` | CPU 架构 |
| `bitness` | `string` | 否 | `{"type":"string"}` | 位数，如 "64" |
| `model` | `string` | 否 | `{"type":"string"}` | 设备型号（SUA） |
| `source` | `optional int32` | 否 | `{"type":"integer"}` | SUA 来源。取值 UserAgentSource：0=未知，1=低熵 CH，2=高熵 CH，3=UA 字符串解析 |
| `ext` | `string` | 否 | `{"$ref":"#/$defs/Ext"}` | 扩展：JSON 对象字符串。 |

## BrandVersion

| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|---|
| `brand` | `string` | 否 | `{"type":"string"}` | 品牌名，如 "Chrome" |
| `version` | `repeated string` | 否 | `{"type":"array","items":{"type":"string"}}` | 版本分量，如 ["192","0","0","0"] |
| `ext` | `string` | 否 | `{"$ref":"#/$defs/Ext"}` | 扩展：JSON 对象字符串。 |

## Geo

| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|---|
| `lat` | `optional double` | 否 | `{"type":"number","minimum":-90,"maximum":90}` | 纬度 |
| `lon` | `optional double` | 否 | `{"type":"number","minimum":-180,"maximum":180}` | 经度 |
| `type` | `optional int32` | 否 | `{"type":"integer"}` | 坐标来源。取值 LocationType：1=GPS，2=IP，3=用户提供 |
| `accuracy` | `optional int32` | 否 | `{"type":"integer","minimum":0}` | 精度（米） |
| `lastfix` | `optional int32` | 否 | `{"type":"integer","minimum":0}` | 坐标年龄（秒） |
| `ipservice` | `optional int32` | 否 | `{"type":"integer"}` | IP 地理库。取值 IpLocationService：1=ip2location，2=Neustar，3=MaxMind，4=NetAcuity |
| `country` | `string` | 否 | `{"type":"string"}` | ISO-3166-1 Alpha-3 等约定国家码 |
| `region` | `string` | 否 | `{"type":"string"}` | 省/州 |
| `metro` | `string` | 否 | `{"type":"string"}` | 都市圈码 |
| `city` | `string` | 否 | `{"type":"string"}` | 城市名。场景：城市级定向。 |
| `zip` | `string` | 否 | `{"type":"string"}` | 邮编 |
| `utcoffset` | `optional int32` | 否 | `{"type":"integer"}` | 相对 UTC 的分钟偏移 |
| `regionfips104` | `string` | 否 | `{"type":"string"}` | FIPS 10-4 地区码（若适用） |
| `ext` | `string` | 否 | `{"$ref":"#/$defs/Ext"}` | 扩展：JSON 对象字符串。 |

## User

| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|---|
| `id` | `string` | 否 | `{"type":"string"}` | Exchange 侧用户 ID |
| `buyeruid` | `string` | 否 | `{"type":"string"}` | 该买家 cookie/映射 ID；DSP 频控关键 |
| `keywords` | `string` | 否 | `{"type":"string"}` | 旧版逗号关键词 |
| `kwarray` | `repeated string` | 否 | `{"type":"array","items":{"type":"string"}}` | 兴趣/上下文关键词 |
| `customdata` | `string` | 否 | `{"type":"string"}` | Exchange 回传的买家自定义串 |
| `geo` | `Geo` | 否 | `{"$ref":"#/$defs/Geo"}` | 用户常住/注册地（可与 Device.geo 不同） |
| `data` | `repeated Data` | 否 | `{"type":"array","items":{"$ref":"#/$defs/Data"}}` | 人群/第一方数据分段 |
| `consent` | `string` | 否 | `{"type":"string"}` | GDPR TCF 同意串等 |
| `eids` | `repeated EID` | 否 | `{"type":"array","items":{"$ref":"#/$defs/EID"}}` | 扩展/统一身份（UID2、Living Docs 等） |
| `ext` | `string` | 否 | `{"$ref":"#/$defs/Ext"}` | 扩展：JSON 对象字符串。 |

## Data

| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|---|
| `id` | `string` | 否 | `{"type":"string"}` | 数据提供方 ID。 |
| `name` | `string` | 否 | `{"type":"string"}` | 数据源名称 |
| `cids` | `repeated string` | 否 | `{"type":"array","items":{"type":"string"}}` | 相关活动/客户 ID |
| `segment` | `repeated Segment` | 否 | `{"type":"array","items":{"$ref":"#/$defs/Segment"}}` | 具体分群 |
| `ext` | `string` | 否 | `{"$ref":"#/$defs/Ext"}` | 扩展：JSON 对象字符串。 |

## Segment

| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|---|
| `id` | `string` | 否 | `{"type":"string"}` | 分群 ID |
| `name` | `string` | 否 | `{"type":"string"}` | 分群名称。 |
| `value` | `string` | 否 | `{"type":"string"}` | 分群取值 |
| `ext` | `string` | 否 | `{"$ref":"#/$defs/Ext"}` | 扩展：JSON 对象字符串。 |

## EID

| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|---|
| `source` | `string` | 否 | `{"type":"string"}` | ID 来源域或体系，如 uidapi.com |
| `uids` | `repeated UID` | 否 | `{"type":"array","items":{"$ref":"#/$defs/UID"}}` | 具体 ID 值列表 |
| `inserter` | `string` | 否 | `{"type":"string"}` | 写入方 |
| `matcher` | `string` | 否 | `{"type":"string"}` | 匹配方 |
| `mm` | `optional int32` | 否 | `{"type":"integer"}` | 匹配方法。取值 IdMatchMethod：0=未知，1=未匹配（直接 cookie/IFA），2=浏览器 cookie sync，3=认证匹配，4=第一方观测，5=推断；>=500 厂商自定义。 |
| `ext` | `string` | 否 | `{"$ref":"#/$defs/Ext"}` | 扩展：JSON 对象字符串。 |

## UID

| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|---|
| `id` | `string` | 否 | `{"type":"string"}` | 用户标识值 |
| `atype` | `optional int32` | 否 | `{"type":"integer"}` | agent 类型。取值 AgentType：1=浏览器/设备，2=App 内，3=跨设备人物；>=500 厂商 |
| `ext` | `string` | 否 | `{"$ref":"#/$defs/Ext"}` | 扩展：JSON 对象字符串。 |

## Source

| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|---|
| `fd` | `optional int32` | 否 | `{"type":"integer","enum":[0,1]}` | 最终决策是否在上游。取值 FlagBool：0=本 Exchange 决策，1=上游最终决策 |
| `tid` | `string` | 否 | `{"type":"string"}` | 事务 ID；链路追踪 |
| `pchain` | `string` | 否 | `{"type":"string"}` | 支付 ID 链（旧字段，优先 schain） |
| `schain` | `SupplyChain` | 否 | `{"$ref":"#/$defs/SupplyChain"}` | SupplyChain 对象 |
| `ext` | `string` | 否 | `{"$ref":"#/$defs/Ext"}` | 扩展：JSON 对象字符串。 |

## SupplyChain

| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|---|
| `complete` | `optional int32` | 是 | `{"type":"integer","enum":[0,1]}` | 链路是否完整。取值 FlagBool：0=可能截断，1=完整 |
| `nodes` | `repeated SupplyChainNode` | 是 | `{"type":"array","minItems":1,"items":{"$ref":"#/$defs/SupplyChainNode"}}` | 从源头到当前的转售节点，按序 |
| `ver` | `string` | 是 | `{"$ref":"#/$defs/NonEmptyString"}` | schain 版本，如 "1.0" |
| `ext` | `string` | 否 | `{"$ref":"#/$defs/Ext"}` | 扩展：JSON 对象字符串。 |

## SupplyChainNode

| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|---|
| `asi` | `string` | 是 | `{"$ref":"#/$defs/NonEmptyString"}` | 广告系统标识域（Advertising System Identifier） |
| `sid` | `string` | 是 | `{"$ref":"#/$defs/NonEmptyString"}` | 该系统上的卖家 ID |
| `rid` | `string` | 否 | `{"type":"string"}` | 该节点请求 ID |
| `name` | `string` | 否 | `{"type":"string"}` | 组织名（可选） |
| `domain` | `string` | 否 | `{"type":"string"}` | 业务域（可选） |
| `hp` | `optional int32` | 是 | `{"type":"integer","enum":[0,1]}` | 是否参与支付路径。取值 FlagBool：0=否，1=是（helps payment） |
| `ext` | `string` | 否 | `{"$ref":"#/$defs/Ext"}` | 扩展：JSON 对象字符串。 |

## Regs

| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|---|
| `coppa` | `optional int32` | 否 | `{"type":"integer","enum":[0,1]}` | 是否适用 COPPA。取值 FlagBool：0=否，1=儿童流量适用 |
| `gdpr` | `optional int32` | 否 | `{"type":"integer","enum":[0,1]}` | 是否适用 GDPR。取值 FlagBool：0=否，1=适用 |
| `us_privacy` | `string` | 否 | `{"type":"string"}` | CCPA/US Privacy 字符串 |
| `gpp` | `string` | 否 | `{"type":"string"}` | Global Privacy Platform 串 |
| `gpp_sid` | `repeated int32` | 否 | `{"type":"array","items":{"type":"integer"}}` | GPP 区段 ID 列表 |
| `ext` | `string` | 否 | `{"$ref":"#/$defs/Ext"}` | 扩展：JSON 对象字符串。 |

## BidResponse

| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|---|
| `id` | `string` | 是 | `{"$ref":"#/$defs/NonEmptyString"}` | 必须等于 BidRequest.id。必填。 |
| `seatbid` | `repeated SeatBid` | 否 | `{"type":"array","items":{"$ref":"#/$defs/SeatBid"}}` | 席位出价集合。真正出价时至少 1 个 SeatBid。 |
| `bidid` | `string` | 否 | `{"type":"string"}` | Bidder 生成的响应 ID。场景：Bidder 侧日志与对账。 |
| `cur` | `string` | 是 | `{"$ref":"#/$defs/Iso4217"}` | 本响应出价币种（ISO-4217）。必填。构建器默认 USD。须落在请求 cur 允许集合内。 |
| `customdata` | `string` | 否 | `{"type":"string"}` | 回传给 Exchange、之后可能再塞进 User.customdata 的不透明串。 |
| `nbr` | `optional int32` | 否 | `{"type":"integer"}` | 不竞价原因。取值 NoBidReason：0=未知错误，1=技术错误，2=非法请求，3=爬虫，4=非人，5=云/代理IP，6=不支持设备，7=屏蔽发布商，8=用户未匹配，9=日读者上限，10=日域名上限；>=500 自定义。 |
| `ext` | `string` | 否 | `{"$ref":"#/$defs/Ext"}` | 扩展：JSON 对象字符串。 |

## SeatBid

| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|---|
| `bid` | `repeated Bid` | 是 | `{"type":"array","minItems":1,"items":{"$ref":"#/$defs/Bid"}}` | 该 seat 的出价列表。至少 1 条才有意义。 |
| `seat` | `string` | 否 | `{"type":"string"}` | 买家席位名。场景：同一 DSP 多 seat 时报价与结算区分。 |
| `group` | `optional int32` | 否 | `{"type":"integer","enum":[0,1]}` | 是否打包整组。取值 FlagBool：0=可部分赢，1=组内须全赢。场景：多 Imp 打包售卖。 |
| `ext` | `string` | 否 | `{"$ref":"#/$defs/Ext"}` | 扩展：JSON 对象字符串。 |

## Bid

| 字段 | protobuf 类型 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|---|
| `id` | `string` | 是 | `{"$ref":"#/$defs/NonEmptyString"}` | Bidder 侧出价 ID。必填。 |
| `impid` | `string` | 是 | `{"$ref":"#/$defs/NonEmptyString"}` | 对应 Imp.id。必填。对不上则 Exchange 丢弃。 |
| `price` | `optional double` | 是 | `{"type":"number","exclusiveMinimum":0}` | CPM 出价，必须 > 0。必填。币种见 BidResponse.cur。 |
| `nurl` | `string` | 否 | `{"type":"string"}` | 赢价通知 URL（可含 ${AUCTION_PRICE}）。场景：胜出后 Exchange 服务端调用。 |
| `burl` | `string` | 否 | `{"type":"string"}` | 计费通知 URL。场景：实际计费时点回调（可与展示分离）。 |
| `lurl` | `string` | 否 | `{"type":"string"}` | 丢单通知 URL。场景：未胜出时的反馈学习。 |
| `adm` | `string` | 否 | `{"type":"string"}` | 广告 markup（HTML/VAST/Native JSON 等）。强烈推荐内联；也可靠 nurl 拉取。 |
| `adid` | `string` | 否 | `{"type":"string"}` | 广告 ID（广告主侧）。场景：广告维度报表。 |
| `adomain` | `repeated string` | 否 | `{"type":"array","items":{"type":"string"}}` | 广告主域名列表。强烈推荐。场景：badv 屏蔽检查与品牌披露。 |
| `bundle` | `string` | 否 | `{"type":"string"}` | 推广应用 bundle。场景：App 安装类广告。 |
| `iurl` | `string` | 否 | `{"type":"string"}` | 抽样创意预览图 URL。场景：审核与质检。 |
| `cid` | `string` | 否 | `{"type":"string"}` | 活动 ID。 |
| `crid` | `string` | 否 | `{"type":"string"}` | 创意 ID。强烈推荐。场景：创意审核、频控、拒登排查。 |
| `tactic` | `string` | 否 | `{"type":"string"}` | 策略战术 ID。场景：DSP 内部策略归因。 |
| `cattax` | `optional int32` | 否 | `{"type":"integer"}` | 类目 taxonomy。取值 CategoryTaxonomy（见枚举）。 |
| `cat` | `repeated string` | 否 | `{"type":"array","items":{"type":"string"}}` | 创意/广告主类目。 |
| `attr` | `repeated int32` | 否 | `{"type":"array","items":{"type":"integer"}}` | 创意属性。取值 CreativeAttribute（见枚举）；>=500 厂商。 |
| `apis` | `repeated int32` | 否 | `{"type":"array","items":{"type":"integer"}}` | 创意使用的 API。取值 ApiFramework（见枚举）。 |
| `protocol` | `optional int32` | 否 | `{"type":"integer"}` | 视频协议。取值 Protocol（见枚举）。 |
| `qagmediarating` | `optional int32` | 否 | `{"type":"integer"}` | 媒体评级。取值 MediaRating：1=全年龄，2=12+，3=成人。 |
| `language` | `string` | 否 | `{"type":"string"}` | 创意语言 ISO-639-1。 |
| `langb` | `string` | 否 | `{"type":"string"}` | 创意语言 BCP-47。 |
| `dealid` | `string` | 否 | `{"type":"string"}` | 成交的 Deal.id。走 PMP 时必填对应 deal。 |
| `w` | `optional int32` | 否 | `{"type":"integer","minimum":1}` | 素材展示宽度，单位为设备无关像素（DIPs）。 |
| `h` | `optional int32` | 否 | `{"type":"integer","minimum":1}` | 素材展示高度，单位为设备无关像素（DIPs）。 |
| `wratio` | `optional int32` | 否 | `{"type":"integer"}` | 创意宽比（原生/弹性）。 |
| `hratio` | `optional int32` | 否 | `{"type":"integer"}` | 创意高比。 |
| `exp` | `optional int32` | 否 | `{"type":"integer"}` | 该出价建议的过期秒数。 |
| `dur` | `optional int32` | 否 | `{"type":"integer"}` | 视频/音频时长秒数。 |
| `mtype` | `MarkupType` | 否 | `{"type":"integer","enum":[1,2,3,4]}` | 素材形态（MarkupType）：1=Banner，2=Video，3=Audio，4=Native。对应 Imp 提供多种形态时，BidCheck 要求显式指定。JSON Schema 不允许 0；protobuf 的 0 只表示未设置。 |
| `slotinpod` | `optional int32` | 否 | `{"type":"integer"}` | 声明占用的 pod 槽位。取值 SlotInPod；特殊 -1=最后槽。场景：CTV pod 选槽出价。 |
| `ext` | `string` | 否 | `{"$ref":"#/$defs/Ext"}` | 扩展：JSON 对象字符串。 |

## Schema 基础类型

```json
{
  "Ext": {
    "title": "Extension object",
    "description": "Exchange- or bidder-specific extensions. Any object may include ext.",
    "type": "object",
    "additionalProperties": true
  },
  "NonEmptyString": {
    "type": "string",
    "minLength": 1
  },
  "Iso4217": {
    "type": "string",
    "pattern": "^[A-Z]{3}$"
  },
  "AuctionType": {
    "type": "integer",
    "description": "1 = first price, 2 = second price plus, >=500 exchange-specific. Deal.at may also be 3 (fixed price)."
  }
}
```

## Native.request 内嵌 JSON

以下为本仓库 Native 请求 Schema 的支持范围，不包含 Native 响应完整协议。先解析外层 request 字符串，再按 native.schema.json 校验；外层字符串类型检查不等于内层合同校验。

枚举语义参考 [IAB Native 1.2](https://github.com/InteractiveAdvertisingBureau/Native-Ads/blob/main/OpenRTB-Native-Ads-Specification-Final-1.2.md)。当前 Schema 对多数分类字段只验证整数类型，未实现完整标准的枚举和组合约束。

### NativeRequest

| 字段 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|
| `ver` | 否 | `{"type":"string"}` | Native 请求协议版本，例如 1.2。 |
| `context` | 否 | `{"type":"integer"}` | 广告所处内容环境的分类。 |
| `contextsubtype` | 否 | `{"type":"integer"}` | 内容环境的细分类。 |
| `plcmttype` | 否 | `{"type":"integer"}` | Native 版位类型：1=信息流，2=原子内容单元，3=内容外部，4=推荐组件。 |
| `plcmtcnt` | 否 | `{"type":"integer"}` | 本次版位可容纳的相同广告单元数量。 |
| `seq` | 否 | `{"type":"integer"}` | 同一上下文中的广告序号。 |
| `assets` | 是 | `{"type":"array","minItems":1,"items":{"$ref":"#/$defs/Asset"}}` | 请求的素材资源列表，至少一个资源。 |
| `aurlsupport` | 否 | `{"type":"integer","enum":[0,1]}` | 是否支持通过资源 URL 返回素材：0=否，1=是。 |
| `durlsupport` | 否 | `{"type":"integer","enum":[0,1]}` | 是否支持动态内容 URL：0=否，1=是。 |
| `eventtrackers` | 否 | `{"type":"array","items":{"type":"object","additionalProperties":true,"required":["event","methods"],"properties":{"event":{"type":"integer"},"methods":{"type":"array","items":{"type":"integer"}},"ext":{"type":"object"}}}}` | 支持的事件及其追踪方式。 |
| `privacy` | 否 | `{"type":"integer","enum":[0,1]}` | 是否支持买方提供隐私声明 URL：0=否，1=是。 |
| `ext` | 否 | `{"type":"object"}` | 交易平台或买方自定义扩展对象。 |

### NativeRequest.eventtrackers[]

| 字段 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|
| `event` | 是 | `{"type":"integer"}` | 待追踪事件类型，例如 1=曝光。 |
| `methods` | 是 | `{"type":"array","items":{"type":"integer"}}` | 支持的追踪方式列表，例如 1=图片像素，2=JavaScript。 |
| `ext` | 否 | `{"type":"object"}` | 交易平台或买方自定义扩展对象。 |

### NativeAsset

| 字段 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|
| `id` | 是 | `{"type":"integer"}` | 请求资源的标识，响应资源使用相同 ID 关联。 |
| `required` | 否 | `{"type":"integer","enum":[0,1]}` | 该资源是否必需：0=可选，1=必需。 |
| `title` | 否 | `{"type":"object","additionalProperties":true,"required":["len"]}` | 标题资源要求。 |
| `img` | 否 | `{"type":"object","additionalProperties":true}` | 图片资源要求。 |
| `video` | 否 | `{"type":"object","additionalProperties":true,"required":["mimes","minduration","maxduration","protocols"]}` | 视频资源要求。 |
| `data` | 否 | `{"type":"object","additionalProperties":true,"required":["type"]}` | 其他数据资源要求，例如描述、品牌或按钮文字。 |
| `ext` | 否 | `{"type":"object"}` | 交易平台或买方自定义扩展对象。 |

### NativeAsset.title

| 字段 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|
| `len` | 是 | `{"type":"integer","minimum":1}` | 标题允许的最大字符数，至少为 1。 |
| `ext` | 否 | `{"type":"object"}` | 交易平台或买方自定义扩展对象。 |

### NativeAsset.img

| 字段 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|
| `type` | 否 | `{"type":"integer"}` | 图片资源类型，例如 1=图标，3=主图。 |
| `w` | 否 | `{"type":"integer"}` | 所需图片的精确宽度，单位为图片像素。 |
| `h` | 否 | `{"type":"integer"}` | 所需图片的精确高度，单位为图片像素。 |
| `wmin` | 否 | `{"type":"integer"}` | 图片最小宽度，单位为图片像素。 |
| `hmin` | 否 | `{"type":"integer"}` | 图片最小高度，单位为图片像素。 |
| `mimes` | 否 | `{"type":"array","items":{"type":"string"}}` | 支持的图片 MIME 类型。 |
| `ext` | 否 | `{"type":"object"}` | 交易平台或买方自定义扩展对象。 |

### NativeAsset.video

| 字段 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|
| `mimes` | 是 | `{"type":"array","items":{"type":"string"}}` | 支持的视频 MIME 类型。 |
| `minduration` | 是 | `{"type":"integer"}` | 最短视频时长，单位为秒。 |
| `maxduration` | 是 | `{"type":"integer"}` | 最长视频时长，单位为秒。 |
| `protocols` | 是 | `{"type":"array","items":{"type":"integer"}}` | 支持的视频响应协议编号。 |
| `ext` | 否 | `{"type":"object"}` | 交易平台或买方自定义扩展对象。 |

### NativeAsset.data

| 字段 | JSON 必填 | JSON 约束 | 中文说明 |
|---|---|---|---|
| `type` | 是 | `{"type":"integer"}` | 数据资源类型，例如 1=赞助方名称，2=描述，12=行动按钮文字。 |
| `len` | 否 | `{"type":"integer"}` | 文本数据允许的最大字符数。 |
| `ext` | 否 | `{"type":"object"}` | 交易平台或买方自定义扩展对象。 |

## protobuf 枚举参考

以下保留源码的英文枚举说明；0 值可能只是 protobuf 未设置哨兵，是否可用于 JSON 取决于字段约束。

### MarkupType

| 名称 | 值 | 说明 |
|---|---|---|
| `MARKUP_TYPE_UNSPECIFIED` | 0 | Unspecified (protobuf default); do not use in the JSON wire format. |
| `MARKUP_TYPE_BANNER` | 1 | Banner/display; corresponds to Imp.banner. |
| `MARKUP_TYPE_VIDEO` | 2 | Video; corresponds to Imp.video. |
| `MARKUP_TYPE_AUDIO` | 3 | Audio; corresponds to Imp.audio. |
| `MARKUP_TYPE_NATIVE` | 4 | Native; corresponds to Imp.native. |

### AuctionType

| 名称 | 值 | 说明 |
|---|---|---|
| `AUCTION_TYPE_UNSPECIFIED` | 0 | Unspecified. |
| `AUCTION_TYPE_FIRST_PRICE` | 1 | First price: the winner pays its own bid. |
| `AUCTION_TYPE_SECOND_PRICE_PLUS` | 2 | Second price plus: pay the second-highest price, optionally with an increment; commonly used for the now-required BidRequest.at. |
| `AUCTION_TYPE_FIXED_PRICE` | 3 | Fixed price; Deal.at only. |

### BannerAdType

| 名称 | 值 | 说明 |
|---|---|---|
| `BANNER_AD_TYPE_UNSPECIFIED` | 0 | Unspecified. |
| `BANNER_AD_TYPE_XHTML_TEXT` | 1 | XHTML text ad, typically mobile. |
| `BANNER_AD_TYPE_XHTML_BANNER` | 2 | XHTML Banner |
| `BANNER_AD_TYPE_JAVASCRIPT` | 3 | JavaScript ad. |
| `BANNER_AD_TYPE_IFRAME` | 4 | iframe |

### CreativeAttribute

| 名称 | 值 | 说明 |
|---|---|---|
| `CREATIVE_ATTRIBUTE_UNSPECIFIED` | 0 | Unspecified. |
| `CREATIVE_ATTRIBUTE_AUDIO_AUTOPLAY` | 1 | Audio ad with autoplay. |
| `CREATIVE_ATTRIBUTE_AUDIO_USER_INITIATED` | 2 | User-initiated audio ad. |
| `CREATIVE_ATTRIBUTE_EXPANDABLE_AUTOMATIC` | 3 | Automatically expandable. |
| `CREATIVE_ATTRIBUTE_EXPANDABLE_CLICK` | 4 | Expandable on click. |
| `CREATIVE_ATTRIBUTE_EXPANDABLE_ROLLOVER` | 5 | Expandable on rollover. |
| `CREATIVE_ATTRIBUTE_IN_BANNER_VIDEO_AUTOPLAY` | 6 | Autoplay in-banner video. |
| `CREATIVE_ATTRIBUTE_IN_BANNER_VIDEO_USER` | 7 | User-initiated in-banner video. |
| `CREATIVE_ATTRIBUTE_POP` | 8 | Pop-over, pop-under or pop-on-exit. |
| `CREATIVE_ATTRIBUTE_PROVOCATIVE` | 9 | Provocative or suggestive imagery. |
| `CREATIVE_ATTRIBUTE_ANNOYING` | 10 | Shaking, flashing or extreme animation. |
| `CREATIVE_ATTRIBUTE_SURVEYS` | 11 | Surveys. |
| `CREATIVE_ATTRIBUTE_TEXT_ONLY` | 12 | Text only. |
| `CREATIVE_ATTRIBUTE_USER_INTERACTIVE` | 13 | User interaction, such as an embedded game. |
| `CREATIVE_ATTRIBUTE_ALERT_STYLE` | 14 | System dialog or alert styling. |
| `CREATIVE_ATTRIBUTE_HAS_AUDIO_TOGGLE` | 15 | Includes an audio toggle. |
| `CREATIVE_ATTRIBUTE_HAS_SKIP_BUTTON` | 16 | Includes a skip button. |
| `CREATIVE_ATTRIBUTE_FLASH` | 17 | Adobe Flash |
| `CREATIVE_ATTRIBUTE_RESPONSIVE` | 18 | Responsive, fluid or without fixed dimensions. |

### AdPosition

| 名称 | 值 | 说明 |
|---|---|---|
| `AD_POSITION_UNKNOWN` | 0 | Unknown. |
| `AD_POSITION_ABOVE_THE_FOLD` | 1 | Above the fold. |
| `AD_POSITION_LOCKED` | 2 | Locked position. |
| `AD_POSITION_BELOW_THE_FOLD` | 3 | Below the fold. |
| `AD_POSITION_HEADER` | 4 | Header. |
| `AD_POSITION_FOOTER` | 5 | Footer. |
| `AD_POSITION_SIDEBAR` | 6 | Sidebar. |
| `AD_POSITION_FULLSCREEN` | 7 | Fullscreen. |

### ExpandableDirection

| 名称 | 值 | 说明 |
|---|---|---|
| `EXPANDABLE_DIRECTION_UNSPECIFIED` | 0 | Unspecified. |
| `EXPANDABLE_DIRECTION_LEFT` | 1 | Left. |
| `EXPANDABLE_DIRECTION_RIGHT` | 2 | Right. |
| `EXPANDABLE_DIRECTION_UP` | 3 | Up. |
| `EXPANDABLE_DIRECTION_DOWN` | 4 | Down. |
| `EXPANDABLE_DIRECTION_FULLSCREEN` | 5 | Fullscreen. |
| `EXPANDABLE_DIRECTION_RESIZE` | 6 | Resize or minimize. |

### ApiFramework

| 名称 | 值 | 说明 |
|---|---|---|
| `API_FRAMEWORK_UNSPECIFIED` | 0 | Unspecified. |
| `API_FRAMEWORK_VPAID_1` | 1 | VPAID 1.0 |
| `API_FRAMEWORK_VPAID_2` | 2 | VPAID 2.0 |
| `API_FRAMEWORK_MRAID_1` | 3 | MRAID 1.0 |
| `API_FRAMEWORK_ORMMA` | 4 | ORMMA |
| `API_FRAMEWORK_MRAID_2` | 5 | MRAID 2.0 |
| `API_FRAMEWORK_MRAID_3` | 6 | MRAID 3.0 |
| `API_FRAMEWORK_OMID_1` | 7 | OMID 1.0 (Open Measurement) |
| `API_FRAMEWORK_SIMID_1` | 8 | SIMID 1.0 |
| `API_FRAMEWORK_SIMID_1_1` | 9 | SIMID 1.1 |

### Protocol

| 名称 | 值 | 说明 |
|---|---|---|
| `PROTOCOL_UNSPECIFIED` | 0 | Unspecified. |
| `PROTOCOL_VAST_1_0` | 1 | VAST 1.0 |
| `PROTOCOL_VAST_2_0` | 2 | VAST 2.0 |
| `PROTOCOL_VAST_3_0` | 3 | VAST 3.0 |
| `PROTOCOL_VAST_1_0_WRAPPER` | 4 | VAST 1.0 Wrapper |
| `PROTOCOL_VAST_2_0_WRAPPER` | 5 | VAST 2.0 Wrapper |
| `PROTOCOL_VAST_3_0_WRAPPER` | 6 | VAST 3.0 Wrapper |
| `PROTOCOL_VAST_4_0` | 7 | VAST 4.0 |
| `PROTOCOL_VAST_4_0_WRAPPER` | 8 | VAST 4.0 Wrapper |
| `PROTOCOL_DAAST_1_0` | 9 | DAAST 1.0 |
| `PROTOCOL_DAAST_1_0_WRAPPER` | 10 | DAAST 1.0 Wrapper |
| `PROTOCOL_VAST_4_1` | 11 | VAST 4.1 |
| `PROTOCOL_VAST_4_1_WRAPPER` | 12 | VAST 4.1 Wrapper |
| `PROTOCOL_VAST_4_2` | 13 | VAST 4.2 |
| `PROTOCOL_VAST_4_2_WRAPPER` | 14 | VAST 4.2 Wrapper |
| `PROTOCOL_VAST_4_3` | 15 | VAST 4.3 |
| `PROTOCOL_VAST_4_3_WRAPPER` | 16 | VAST 4.3 Wrapper |

### VideoPlcmt

| 名称 | 值 | 说明 |
|---|---|---|
| `VIDEO_PLCMT_UNSPECIFIED` | 0 | Unspecified. |
| `VIDEO_PLCMT_INSTREAM` | 1 | Instream with main content, sound on by default; pre-, mid- or post-roll. |
| `VIDEO_PLCMT_ACCOMPANYING_CONTENT` | 2 | Accompanying text/image content; playback starts in the viewport. |
| `VIDEO_PLCMT_INTERSTITIAL` | 3 | Interstitial filling the viewport and not scrollable out of view. |
| `VIDEO_PLCMT_NO_CONTENT` | 4 | No video content: slideshows, feeds or floating players. |
| `VIDEO_PLCMT_PAUSE` | 5 | Shown when the user pauses content (new AdCOM value). |
| `VIDEO_PLCMT_SCREENSAVER` | 6 | Screensaver placement. |
| `VIDEO_PLCMT_OVERLAY` | 7 | Overlay on content, outside a traditional ad break. |
| `VIDEO_PLCMT_SQUEEZEBACK` | 8 | Content shrinks to share the screen with an ad. |
| `VIDEO_PLCMT_IN_SCENE` | 9 | Placement within the content scene. |

### VideoLinearity

| 名称 | 值 | 说明 |
|---|---|---|
| `VIDEO_LINEARITY_UNSPECIFIED` | 0 | Unspecified. |
| `VIDEO_LINEARITY_LINEAR` | 1 | Linear ad, including video assets. |
| `VIDEO_LINEARITY_NON_LINEAR` | 2 | Non-linear/overlay. |

### PlaybackMethod

| 名称 | 值 | 说明 |
|---|---|---|
| `PLAYBACK_METHOD_UNSPECIFIED` | 0 | Unspecified. |
| `PLAYBACK_METHOD_PAGE_LOAD_SOUND_ON` | 1 | Page-load playback with sound on. |
| `PLAYBACK_METHOD_PAGE_LOAD_SOUND_OFF` | 2 | Page-load playback muted by default. |
| `PLAYBACK_METHOD_CLICK_SOUND_ON` | 3 | Click-to-play with sound on. |
| `PLAYBACK_METHOD_MOUSE_OVER_SOUND_ON` | 4 | Mouse-over playback with sound on. |
| `PLAYBACK_METHOD_VIEWPORT_SOUND_ON` | 5 | Viewport-triggered playback with sound on. |
| `PLAYBACK_METHOD_VIEWPORT_SOUND_OFF` | 6 | Viewport-triggered playback muted by default. |
| `PLAYBACK_METHOD_CONTINUOUS` | 7 | Continuous playback until stopped by the user. |

### PlaybackCessationMode

| 名称 | 值 | 说明 |
|---|---|---|
| `PLAYBACK_CESSATION_UNSPECIFIED` | 0 | Unspecified. |
| `PLAYBACK_CESSATION_COMPLETION_OR_USER` | 1 | Stop on completion or user action. |
| `PLAYBACK_CESSATION_LEAVING_VIEWPORT_OR_USER` | 2 | Stop on leaving the viewport or user action. |
| `PLAYBACK_CESSATION_FLOATING_UNTIL_COMPLETION` | 3 | Continue in a floating player after leaving the viewport until completion or user action. |

### DeliveryMethod

| 名称 | 值 | 说明 |
|---|---|---|
| `DELIVERY_METHOD_UNSPECIFIED` | 0 | Unspecified. |
| `DELIVERY_METHOD_STREAMING` | 1 | Streaming. |
| `DELIVERY_METHOD_PROGRESSIVE` | 2 | Progressive download. |
| `DELIVERY_METHOD_DOWNLOAD` | 3 | Full download. |

### CompanionType

| 名称 | 值 | 说明 |
|---|---|---|
| `COMPANION_TYPE_UNSPECIFIED` | 0 | Unspecified. |
| `COMPANION_TYPE_STATIC` | 1 | Static Resource |
| `COMPANION_TYPE_HTML` | 2 | HTML Resource |
| `COMPANION_TYPE_IFRAME` | 3 | iframe Resource |

### VideoCompanionMode

| 名称 | 值 | 说明 |
|---|---|---|
| `VIDEO_COMPANION_MODE_CONCURRENT` | 0 | Shown concurrently with video. |
| `VIDEO_COMPANION_MODE_END_CARD` | 1 | End card after video completion. |

### DeviceType

| 名称 | 值 | 说明 |
|---|---|---|
| `DEVICE_TYPE_UNSPECIFIED` | 0 | Unspecified. |
| `DEVICE_TYPE_MOBILE_TABLET_GENERAL` | 1 | General mobile device when phone/tablet (4/5) is unknown. |
| `DEVICE_TYPE_PERSONAL_COMPUTER` | 2 | Desktop or laptop browser. |
| `DEVICE_TYPE_CONNECTED_TV` | 3 | Smart TV / CTV. |
| `DEVICE_TYPE_PHONE` | 4 | Phone. |
| `DEVICE_TYPE_TABLET` | 5 | Tablet. |
| `DEVICE_TYPE_CONNECTED_DEVICE` | 6 | Non-TV connected device, such as a console or streaming box. |
| `DEVICE_TYPE_SET_TOP_BOX` | 7 | Operator set-top box. |
| `DEVICE_TYPE_OOH_DEVICE` | 8 | Digital out-of-home screen. |

### ConnectionType

| 名称 | 值 | 说明 |
|---|---|---|
| `CONNECTION_TYPE_UNKNOWN` | 0 | Unknown. |
| `CONNECTION_TYPE_ETHERNET` | 1 | Wired Ethernet. |
| `CONNECTION_TYPE_WIFI` | 2 | Wi‑Fi |
| `CONNECTION_TYPE_CELLULAR_UNKNOWN` | 3 | Cellular, generation unknown. |
| `CONNECTION_TYPE_CELLULAR_2G` | 4 | Cellular 2G. |
| `CONNECTION_TYPE_CELLULAR_3G` | 5 | Cellular 3G. |
| `CONNECTION_TYPE_CELLULAR_4G` | 6 | Cellular 4G. |
| `CONNECTION_TYPE_CELLULAR_5G` | 7 | Cellular 5G. |

### LocationType

| 名称 | 值 | 说明 |
|---|---|---|
| `LOCATION_TYPE_UNSPECIFIED` | 0 | Unspecified. |
| `LOCATION_TYPE_GPS` | 1 | GPS or operating-system location services. |
| `LOCATION_TYPE_IP` | 2 | Inferred from IP address. |
| `LOCATION_TYPE_USER_PROVIDED` | 3 | User-provided, such as registration data. |

### IpLocationService

| 名称 | 值 | 说明 |
|---|---|---|
| `IP_LOCATION_SERVICE_UNSPECIFIED` | 0 | Unspecified. |
| `IP_LOCATION_SERVICE_IP2LOCATION` | 1 | IP2Location |
| `IP_LOCATION_SERVICE_NEUSTAR` | 2 | Neustar (Quova) |
| `IP_LOCATION_SERVICE_MAXMIND` | 3 | MaxMind |
| `IP_LOCATION_SERVICE_NETACUITY` | 4 | Digital Element |

### ContentContext

| 名称 | 值 | 说明 |
|---|---|---|
| `CONTENT_CONTEXT_UNSPECIFIED` | 0 | Unspecified. |
| `CONTENT_CONTEXT_VIDEO` | 1 | Video file or stream. |
| `CONTENT_CONTEXT_GAME` | 2 | Game. |
| `CONTENT_CONTEXT_MUSIC` | 3 | Music or radio stream. |
| `CONTENT_CONTEXT_APPLICATION` | 4 | Application. |
| `CONTENT_CONTEXT_TEXT` | 5 | Text page or article. |
| `CONTENT_CONTEXT_OTHER` | 6 | Other. |
| `CONTENT_CONTEXT_UNKNOWN` | 7 | Unknown. |

### ProductionQuality

| 名称 | 值 | 说明 |
|---|---|---|
| `PRODUCTION_QUALITY_UNKNOWN` | 0 | Unknown. |
| `PRODUCTION_QUALITY_PROFESSIONAL` | 1 | Professional production. |
| `PRODUCTION_QUALITY_PROSUMER` | 2 | Prosumer production. |
| `PRODUCTION_QUALITY_USER_GENERATED` | 3 | UGC |

### MediaRating

| 名称 | 值 | 说明 |
|---|---|---|
| `MEDIA_RATING_UNSPECIFIED` | 0 | Unspecified. |
| `MEDIA_RATING_ALL_AUDIENCES` | 1 | All audiences. |
| `MEDIA_RATING_OVER_12` | 2 | Ages 12 and above. |
| `MEDIA_RATING_MATURE` | 3 | Mature audiences. |

### CategoryTaxonomy

| 名称 | 值 | 说明 |
|---|---|---|
| `CATEGORY_TAXONOMY_UNSPECIFIED` | 0 | Unspecified. |
| `CATEGORY_TAXONOMY_IAB_CONTENT_1_0` | 1 | IAB Content Category 1.0 (no longer recommended). |
| `CATEGORY_TAXONOMY_IAB_CONTENT_2_0` | 2 | IAB Content Category 2.0 (no longer recommended). |
| `CATEGORY_TAXONOMY_IAB_AD_PRODUCT_1_0` | 3 | IAB Ad Product Taxonomy 1.0 |
| `CATEGORY_TAXONOMY_IAB_AUDIENCE_1_1` | 4 | IAB Audience Taxonomy 1.1 |
| `CATEGORY_TAXONOMY_IAB_CONTENT_2_1` | 5 | IAB Content Taxonomy 2.1 |
| `CATEGORY_TAXONOMY_IAB_CONTENT_2_2` | 6 | IAB Content Taxonomy 2.2 |
| `CATEGORY_TAXONOMY_IAB_CONTENT_3_0` | 7 | IAB Content Taxonomy 3.0, where supported. |

### FeedType

| 名称 | 值 | 说明 |
|---|---|---|
| `FEED_TYPE_UNSPECIFIED` | 0 | Unspecified. |
| `FEED_TYPE_MUSIC_SERVICE` | 1 | Music service. |
| `FEED_TYPE_FM_AM_BROADCAST` | 2 | FM/AM radio. |
| `FEED_TYPE_PODCAST` | 3 | Podcast. |

### VolumeNormalizationMode

| 名称 | 值 | 说明 |
|---|---|---|
| `VOLUME_NORM_NONE` | 0 | None. |
| `VOLUME_NORM_AVERAGE` | 1 | Normalize to average content volume. |
| `VOLUME_NORM_PEAK` | 2 | Normalize to peak content volume. |
| `VOLUME_NORM_LOUDNESS` | 3 | Loudness normalization. |
| `VOLUME_NORM_CUSTOM` | 4 | Custom normalization. |

### AgentType

| 名称 | 值 | 说明 |
|---|---|---|
| `AGENT_TYPE_UNSPECIFIED` | 0 | Unspecified. |
| `AGENT_TYPE_WEB_OR_DEVICE` | 1 | Browser/device-level ID, such as a cookie. |
| `AGENT_TYPE_IN_APP` | 2 | In-app ID, such as a device advertising ID. |
| `AGENT_TYPE_PERSON` | 3 | Person-level ID across devices. |

### NoBidReason

| 名称 | 值 | 说明 |
|---|---|---|
| `NO_BID_REASON_UNKNOWN_ERROR` | 0 | Unknown error. |
| `NO_BID_REASON_TECHNICAL_ERROR` | 1 | Technical error. |
| `NO_BID_REASON_INVALID_REQUEST` | 2 | Invalid request. |
| `NO_BID_REASON_KNOWN_WEB_SPIDER` | 3 | Known crawler. |
| `NO_BID_REASON_SUSPECTED_NONHUMAN` | 4 | Suspected non-human traffic. |
| `NO_BID_REASON_CLOUD_OR_PROXY_IP` | 5 | Cloud, data-center or proxy IP address. |
| `NO_BID_REASON_UNSUPPORTED_DEVICE` | 6 | Unsupported device. |
| `NO_BID_REASON_BLOCKED_PUBLISHER` | 7 | Blocked publisher or site. |
| `NO_BID_REASON_UNMATCHED_USER` | 8 | Unmatched user. |
| `NO_BID_REASON_DAILY_READER_CAP` | 9 | Daily reader cap. |
| `NO_BID_REASON_DAILY_DOMAIN_CAP` | 10 | Daily domain cap. |

### FlagBool

| 名称 | 值 | 说明 |
|---|---|---|
| `FLAG_FALSE` | 0 | No, off or production, depending on the field. |
| `FLAG_TRUE` | 1 | Yes, on or test, depending on the field. |

### Ssai

| 名称 | 值 | 说明 |
|---|---|---|
| `SSAI_UNKNOWN` | 0 | Unknown. |
| `SSAI_CLIENT` | 1 | Client-side insertion. |
| `SSAI_SERVER` | 2 | Server-side ad insertion (SSAI). |
| `SSAI_MIXED` | 3 | Hybrid client-side/server-side insertion. |

### QtySourceType

| 名称 | 值 | 说明 |
|---|---|---|
| `QTY_SOURCE_MEASUREMENT_VENDOR` | 0 | Measurement-provider common method, such as CMM. |
| `QTY_SOURCE_PUBLISHER` | 1 | Publisher- or vendor-specific method. |

### SlotInPod

| 名称 | 值 | 说明 |
|---|---|---|
| `SLOT_IN_POD_ANY` | 0 | Any slot. |
| `SLOT_IN_POD_FIRST` | 1 | First/start slot only. |
| `SLOT_IN_POD_FIRST_OR_LAST` | 2 | First or last slot. |
| `SLOT_IN_POD_FIRST_MIDDLE_OR_LAST` | 3 | First, middle or last slot. |

### PodSequence

| 名称 | 值 | 说明 |
|---|---|---|
| `POD_SEQUENCE_ANY` | 0 | Any pod. |
| `POD_SEQUENCE_FIRST` | 1 | First pod. |

### SourceRelationship

| 名称 | 值 | 说明 |
|---|---|---|
| `SOURCE_RELATIONSHIP_INDIRECT` | 0 | Indirect. |
| `SOURCE_RELATIONSHIP_DIRECT` | 1 | Direct. |

### UserAgentSource

| 名称 | 值 | 说明 |
|---|---|---|
| `USER_AGENT_SOURCE_UNSPECIFIED` | 0 | Unknown or not applicable. |
| `USER_AGENT_SOURCE_CLIENT_HINTS_LOW` | 1 | Low-entropy Client Hints only. |
| `USER_AGENT_SOURCE_CLIENT_HINTS_HIGH` | 2 | Includes high-entropy Client Hints. |
| `USER_AGENT_SOURCE_USER_AGENT_STRING` | 3 | Parsed from the raw UA string. |

### PodDeduplication

| 名称 | 值 | 说明 |
|---|---|---|
| `POD_DEDUPE_UNSPECIFIED` | 0 | Unspecified protobuf placeholder; the AdCOM list starts at 1. |
| `POD_DEDUPE_ADOMAIN` | 1 | Deduplicate by advertiser domain (adomain). |
| `POD_DEDUPE_IAB_CATEGORY` | 2 | Deduplicate by IAB Tech Lab Content Taxonomy category. |
| `POD_DEDUPE_CREATIVE_ID` | 3 | Deduplicate by creative ID. |
| `POD_DEDUPE_MEDIAFILE_URL` | 4 | Deduplicate by mediafile URL. |

### AutoRefreshTrigger

| 名称 | 值 | 说明 |
|---|---|---|
| `AUTO_REFRESH_TRIGGER_UNKNOWN` | 0 | Unknown. |
| `AUTO_REFRESH_TRIGGER_USER_ACTION` | 1 | Refresh triggered by user action. |
| `AUTO_REFRESH_TRIGGER_EVENT` | 2 | Event-driven refresh, such as content interaction. |
| `AUTO_REFRESH_TRIGGER_TIME` | 3 | Automatic refresh at timed intervals. |

### IdMatchMethod

| 名称 | 值 | 说明 |
|---|---|---|
| `ID_MATCH_METHOD_UNKNOWN` | 0 | Unknown. |
| `ID_MATCH_METHOD_NO_MATCH` | 1 | Unmatched: obtained directly from a third-party cookie or IFA. |
| `ID_MATCH_METHOD_BROWSER_COOKIE_SYNC` | 2 | Real-time browser cookie synchronization. |
| `ID_MATCH_METHOD_AUTHENTICATED` | 3 | Authenticated user match, such as email login or hashed PII. |
| `ID_MATCH_METHOD_OBSERVED` | 4 | Unauthenticated first-party observation, such as GUID, SharedID or session data. |
| `ID_MATCH_METHOD_INFERENCE` | 5 | Inferred across browsers/devices, such as IP plus UA. |
