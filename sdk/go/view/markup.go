package view

import openrtb "github.com/oakrtb/oakrtb/sdk/go/oakrtb/v2"

// MarkupMask is a bitmask of markup types on an Imp (banner/video/audio/native).
// Note: this differs from oakrtb.v2.Format (a Banner size entry).
type MarkupMask uint8

const (
	// MarkupNone indicates that no markup type is set.
	MarkupNone MarkupMask = 0
	// MarkupBanner corresponds to imp.banner and Bid.mtype 1.
	MarkupBanner MarkupMask = 1 << 0
	// MarkupVideo corresponds to imp.video and Bid.mtype 2.
	MarkupVideo MarkupMask = 1 << 1
	// MarkupAudio corresponds to imp.audio and Bid.mtype 3.
	MarkupAudio MarkupMask = 1 << 2
	// MarkupNative corresponds to imp.native and Bid.mtype 4.
	MarkupNative MarkupMask = 1 << 3
)

// Has reports whether the bitmask contains the specified flag.
func (f MarkupMask) Has(flag MarkupMask) bool { return f&flag != 0 }

// HasBanner reports whether Banner markup is present.
func (f MarkupMask) HasBanner() bool { return f.Has(MarkupBanner) }

// HasVideo reports whether Video markup is present.
func (f MarkupMask) HasVideo() bool { return f.Has(MarkupVideo) }

// HasAudio reports whether Audio markup is present.
func (f MarkupMask) HasAudio() bool { return f.Has(MarkupAudio) }

// HasNative reports whether Native markup is present.
func (f MarkupMask) HasNative() bool { return f.Has(MarkupNative) }

// Count returns the number of markup types set in the bitmask.
func (f MarkupMask) Count() int {
	n := 0
	for x := f; x != 0; x >>= 1 {
		n += int(x & 1)
	}
	return n
}

// Primary returns the single set markup bit, or MarkupNone if zero or multiple bits are set.
func (f MarkupMask) Primary() MarkupMask {
	if f.Count() == 1 {
		return f
	}
	return MarkupNone
}

// Mtype maps a single markup bit to Bid.mtype (1–4); returns 0 for zero or multiple bits.
func (f MarkupMask) Mtype() int32 {
	switch f.Primary() {
	case MarkupBanner:
		return 1
	case MarkupVideo:
		return 2
	case MarkupAudio:
		return 3
	case MarkupNative:
		return 4
	default:
		return 0
	}
}

// Inventory is the mutually exclusive BidRequest inventory type (site/app/dooh).
// Note: this differs from proto Content.Channel (content distribution channel).
type Inventory uint8

const (
	// InventoryNone indicates that site/app/dooh is unset.
	InventoryNone Inventory = iota
	// InventorySite indicates website inventory (BidRequest.site).
	InventorySite
	// InventoryApp indicates app inventory (BidRequest.app).
	InventoryApp
	// InventoryDooh indicates digital out-of-home inventory (BidRequest.dooh).
	InventoryDooh
)

// String returns the short inventory type name (site/app/dooh/none).
func (c Inventory) String() string {
	switch c {
	case InventorySite:
		return "site"
	case InventoryApp:
		return "app"
	case InventoryDooh:
		return "dooh"
	default:
		return "none"
	}
}

// MarkupFromImp returns the Imp markup bitmask without full validation.
func MarkupFromImp(imp *openrtb.Imp) MarkupMask {
	if imp == nil {
		return MarkupNone
	}
	var f MarkupMask
	if imp.Banner != nil {
		f |= MarkupBanner
	}
	if imp.Video != nil {
		f |= MarkupVideo
	}
	if imp.Audio != nil {
		f |= MarkupAudio
	}
	if imp.Native != nil {
		f |= MarkupNative
	}
	return f
}

// MarkupFromMtype maps Bid.mtype / MarkupType to a MarkupMask bit (or MarkupNone).
func MarkupFromMtype(m openrtb.MarkupType) MarkupMask {
	switch m {
	case openrtb.MarkupType_MARKUP_TYPE_BANNER:
		return MarkupBanner
	case openrtb.MarkupType_MARKUP_TYPE_VIDEO:
		return MarkupVideo
	case openrtb.MarkupType_MARKUP_TYPE_AUDIO:
		return MarkupAudio
	case openrtb.MarkupType_MARKUP_TYPE_NATIVE:
		return MarkupNative
	default:
		return MarkupNone
	}
}
