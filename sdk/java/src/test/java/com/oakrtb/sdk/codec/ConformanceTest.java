package com.oakrtb.sdk.codec;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.oakrtb.openrtb.v2.*;
import com.oakrtb.sdk.bidcheck.BidCheck;
import com.oakrtb.sdk.validation.BasicValidation;
import com.oakrtb.sdk.view.RequestView;
import org.junit.jupiter.api.Test;
import java.nio.file.Path;
import java.util.ArrayList;
import java.util.Collections;
import static org.junit.jupiter.api.Assertions.*;

class ConformanceTest {
  @Test void sharedContract() throws Exception {
    var mapper = new ObjectMapper().enable(com.fasterxml.jackson.databind.DeserializationFeature.USE_BIG_DECIMAL_FOR_FLOATS);
    var groups = mapper.readTree(Path.of("../../testdata/conformance/cases.json").toFile());
    for (var c : groups.get("roundtrip")) {
      com.google.protobuf.Message decoded;
      if (c.get("kind").asText().equals("request")) {
        var model = Json.parseBidRequest(c.get("input").toString());
        decoded = BidRequest.parseFrom(model.toByteArray());
      } else {
        var model = Json.parseBidResponse(c.get("input").toString());
        decoded = BidResponse.parseFrom(model.toByteArray());
      }
      assertEquals(c.get("expected"), mapper.readTree(Json.toJson(decoded)), c.get("name").asText());
    }
    for (var c : groups.get("reject")) {
      assertThrows(Exception.class, () -> {
        if (c.get("kind").asText().equals("request")) Json.parseBidRequest(c.get("input").toString());
        else Json.parseBidResponse(c.get("input").toString());
      }, c.get("name").asText());
    }
    for (var c : groups.get("basic")) {
      Runnable check;
      if (c.get("kind").asText().equals("request")) {
        var model = Json.parseBidRequest(c.get("input").toString()); check = () -> BasicValidation.validateRequest(model);
      } else {
        var model = Json.parseBidResponse(c.get("input").toString()); check = () -> BasicValidation.validateResponse(model);
      }
      if (c.get("ok").asBoolean()) assertDoesNotThrow(check::run, c.get("name").asText());
      else assertThrows(IllegalArgumentException.class, check::run, c.get("name").asText());
    }
    for (var c : groups.get("bidcheck")) {
      var req = RequestView.of(Json.parseBidRequest(c.get("request").toString()));
      var res = Json.parseBidResponse(c.get("response").toString());
      var result = BidCheck.response(req,res);
      assertEquals(c.get("ok").asBoolean(), result.ok(), c.get("name").asText());
      var codes = result.issues().stream().map(i -> i.code()).sorted().toList();
      var expected = new ArrayList<String>(); for (var code : c.get("codes")) expected.add(code.asText());
      Collections.sort(expected);
      assertEquals(expected,codes,c.get("name").asText());
    }
    for (var c : groups.get("facts")) {
      var v = RequestView.of(Json.parseBidRequest(c.get("request").toString()));
      var f = v.facts().getFirst();
      assertEquals(c.get("banner_w").intValue(), f.bannerW());
      assertNull(f.bannerH());
      assertEquals(c.get("native_request").asText(), f.nativeRequest());
      var r = com.oakrtb.sdk.view.ResponseView.of(Json.parseBidResponse(c.get("response").toString()));
      assertEquals(c.get("has_adm").asBoolean(), r.facts().getFirst().hasAdm());
    }
    for (var c : groups.get("ready")) {
      var req = Json.parseBidRequest(c.get("request").toString());
      var result = com.oakrtb.sdk.validation.Readiness.impReady(req.getImp(0), c.get("mtype").intValue());
      assertEquals(c.get("ok").asBoolean(), result.ok(), c.get("name").asText());
      var codes = new java.util.ArrayList<String>();
      c.get("codes").forEach(code -> codes.add(code.asText()));
      assertEquals(codes.stream().sorted().toList(), result.issues().stream().map(i -> i.code()).sorted().toList());
      boolean builderOk = c.get("builder_ok").asBoolean();
      assertEquals(builderOk, com.oakrtb.sdk.validation.Readiness.impReady(req.getImp(0)).ok());
      org.junit.jupiter.api.function.Executable build = () -> com.oakrtb.sdk.builder.BidRequestBuilder.create(req.getId())
          .auctionType(req.getAt()).currency("USD").addImp(req.getImp(0)).build();
      if (builderOk) assertDoesNotThrow(build); else assertThrows(IllegalStateException.class, build);
    }
  }
}
