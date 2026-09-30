package com.oakrtb.sdk.view;

import com.oakrtb.openrtb.v2.BidRequest;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Optional;

/** Query view of an immutable generated model; no pipeline state. */
public final class RequestView {
    private final BidRequest request;
    private final java.time.Instant deadline;
    private final List<ImpView> imps;

    private RequestView(BidRequest request) {
      this.request = request;
      this.deadline = request.getTmax() > 0 ? java.time.Instant.now().plusMillis(request.getTmax() * 85L / 100L) : null;
      this.imps = request.getImpList().stream().map(ImpView::of).toList();
    }

    /** Validate once and construct a complete query view. */
    public static RequestView of(BidRequest request) {
      com.oakrtb.sdk.validation.BasicValidation.validateRequest(request);
      return new RequestView(request);
    }

    public List<ImpView> imps() { return imps; }
    public java.time.Instant deadline() { return deadline; }

    /**
     * Auction ID.
     *
     * @return the request ID
     */
    public String auctionId() {
      return request.getId();
    }

    /**
     * Auction type (at).
     *
     * @return the at value
     */
    public int auctionType() {
      return request.getAt();
    }

    /**
     * Inventory type ({@link Inventory}, distinct from Content.Channel).
     *
     * @return the Inventory enum
     */
    public Inventory inventory() {
      return request.hasSite() ? Inventory.SITE : request.hasApp() ? Inventory.APP : request.hasDooh() ? Inventory.DOOH : Inventory.NONE;
    }

    /**
     * Allowed bid currencies.
     *
     * @return the currency code list
     */
    public List<String> currencies() {
      return request.getCurList();
    }

    /**
     * Reports whether the suggested deadline has passed.
     *
     * @return true if expired
     */
    public boolean pastDeadline() {
      return deadline != null && java.time.Instant.now().isAfter(deadline);
    }

    /**
     * Finds an Imp view by ID.
     *
     * @param id Imp id
     * @return the matching ImpView, or empty if absent
     */
    public Optional<ImpView> findImp(String id) {
      if (id == null) {
        return Optional.empty();
      }
      return imps.stream().filter(i -> id.equals(i.id())).findFirst();
    }

    /**
     * Filters Imps by the specified format bit (e.g. {@link MarkupMask#BANNER}).
     *
     * @param formatFlag single-bit format constant
     * @return matching Imps
     */
    public List<ImpView> impsWith(int formatFlag) {
      List<ImpView> out = new ArrayList<>();
      for (ImpView iv : imps) {
        if (iv.markup().has(formatFlag)) {
          out.add(iv);
        }
      }
      return Collections.unmodifiableList(out);
    }

    /**
     * Extracts compact per-Imp facts for matching/bidding.
     *
     * @return the ImpFact list
     */
    public List<ImpFact> facts() {
      List<ImpFact> out = new ArrayList<>(imps.size());
      for (ImpView iv : imps) {
        out.add(ImpFact.from(iv));
      }
      return Collections.unmodifiableList(out);
    }

    /**
     * Original BidRequest reference.
     *
     * @return BidRequest
     */
    public BidRequest request() {
      return request;
    }
    public record ImpFact(
      String id,
      MarkupMask markup,
      int mtype,
      String tagId,
      double bidFloor,
      String bidFloorCur,
      int secure,
      int instl,
      int rwdd,
      int ssai,
      Integer bannerW,
      Integer bannerH,
      String nativeRequest) {

    static ImpFact from(ImpView iv) {
      Integer w = null;
      Integer h = null;
      if (iv.banner() != null) {
        if (iv.banner().hasW()) {
          w = iv.banner().getW();
        }
        if (iv.banner().hasH()) {
          h = iv.banner().getH();
        }
      }
      String nativeReq = null;
      if (iv.nativeAd() != null) {
        String r = iv.nativeAd().getRequest();
        if (r != null && !r.isEmpty()) {
          nativeReq = r;
        }
      }
      return new ImpFact(
          iv.id(),
          iv.markup(),
          iv.markup().mtype(),
          iv.tagId(),
          iv.bidFloor(),
          iv.bidFloorCur(),
          iv.secure(),
          iv.instl(),
          iv.rwdd(),
          iv.ssai(),
          w,
          h,
          nativeReq);
    }

    /** Reports whether Banner format is present. */
    public boolean hasBanner() {
      return markup.hasBanner();
    }

    /** Reports whether Video format is present. */
    public boolean hasVideo() {
      return markup.hasVideo();
    }

    /** Reports whether Audio format is present. */
    public boolean hasAudio() {
      return markup.hasAudio();
    }

    /** Reports whether Native format is present. */
    public boolean hasNative() {
      return markup.hasNative();
    }
  }
}
