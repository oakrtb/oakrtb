package com.oakrtb.sdk.builder;
import com.oakrtb.sdk.codec.Json;

import com.google.protobuf.Descriptors;
import com.google.protobuf.DynamicMessage;
import com.google.protobuf.util.JsonFormat;
import com.oakrtb.openrtb.v2.OpenRtbProto;
import org.junit.jupiter.api.Test;
import static org.junit.jupiter.api.Assertions.*;

class PresenceTest {
  @Test void numericZeroRoundTrip() throws Exception {
    for (var message : OpenRtbProto.getDescriptor().getMessageTypes()) {
      for (var field : message.getFields()) {
        if (field.isRepeated() || (field.getType() != Descriptors.FieldDescriptor.Type.INT32
            && field.getType() != Descriptors.FieldDescriptor.Type.DOUBLE)) continue;
        var absent = DynamicMessage.newBuilder(message).build();
        assertEquals("{}", Json.toJson(absent));
        var original = absent.toBuilder().setField(field, field.getDefaultValue()).build();
        var back = DynamicMessage.newBuilder(message);
        JsonFormat.parser().merge(Json.toJson(original), back);
        assertTrue(back.hasField(field), field.getFullName());
      }
    }
  }
  @Test void zeroNoBidAndSwitchToBid() {
    var builder = BidResponseBuilder.create("a").noBid(0);
    assertTrue(Json.toJson(builder.build()).contains("\"nbr\":0"));
    var bid = com.oakrtb.openrtb.v2.Bid.newBuilder().setId("b").setImpid("1").setPrice(1).build();
    assertFalse(builder.addSeatBid("s", bid).build().hasNbr());
  }
}
