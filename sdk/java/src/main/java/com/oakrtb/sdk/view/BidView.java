package com.oakrtb.sdk.view;

import com.oakrtb.openrtb.v2.Bid;
import com.oakrtb.openrtb.v2.MarkupType;
import java.util.List;

/** One bid with its owning seat and format. */
public record BidView(
      Bid bid,
      String id,
      String impid,
      String seat,
      double price,
      MarkupMask markup,
      int mtype,
      String crid,
      String cid,
      String dealid,
      int w,
      int h,
      int dur,
      String adm,
      String nurl,
      String burl,
      String lurl,
      List<String> adomain) {
  public BidView { adomain = List.copyOf(adomain); }

  static BidView of(Bid bid, String seat) {
    MarkupType mt = bid.getMtype();
    MarkupMask markup = MarkupMask.fromMtype(mt);
    return new BidView(
        bid,
        bid.getId(),
        bid.getImpid(),
        seat,
        bid.getPrice(),
        markup,
        bid.getMtypeValue(),
        bid.getCrid(),
        bid.getCid(),
        bid.getDealid(),
        bid.getW(),
        bid.getH(),
        bid.getDur(),
        bid.getAdm(),
        bid.getNurl(),
        bid.getBurl(),
        bid.getLurl(),
        List.copyOf(bid.getAdomainList()));
  }
}
