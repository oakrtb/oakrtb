package com.oakrtb.sdk.view;

import com.oakrtb.openrtb.v2.SeatBid;
import java.util.List;

/** A seat and its immutable bid projections. */
public record SeatBidView(SeatBid seatBid, String seat, int group, List<BidView> bids) {
  public SeatBidView { bids = List.copyOf(bids); }
  static SeatBidView of(SeatBid seat) {
    return new SeatBidView(seat, seat.getSeat(), seat.getGroup(), seat.getBidList().stream().map(b -> BidView.of(b, seat.getSeat())).toList());
  }
}
