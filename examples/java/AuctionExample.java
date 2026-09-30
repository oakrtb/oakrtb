import com.oakrtb.sdk.bidcheck.BidCheck;
import com.oakrtb.sdk.builder.BidRequestBuilder;
import com.oakrtb.sdk.builder.BidResponseBuilder;
import com.oakrtb.sdk.builder.ImpBuilders;
import com.oakrtb.sdk.builder.Parts;
import com.oakrtb.sdk.codec.Json;
import com.oakrtb.sdk.jsonschema.Schema;
import com.oakrtb.sdk.view.RequestView;
import com.oakrtb.sdk.view.ResponseView;
import java.nio.charset.StandardCharsets;

public class AuctionExample {
  // Simulate the DSP boundary; price == 0 means no candidate was found.
  static byte[] respond(String requestJson, double price) throws Exception {
    // Full contract validation is useful at integration boundaries.
    var report = Schema.request(requestJson);
    if (!report.isOk()) throw new IllegalArgumentException(report.toJson());
    var req = Json.parseBidRequest(requestJson);
    var requestView = RequestView.of(req);
    var imp = requestView.findImp("imp-1").orElseThrow();
    if (!imp.markup().hasBanner()) throw new IllegalArgumentException("expected Banner imp-1");

    var response = BidResponseBuilder.create(requestView.auctionId()).currency("USD");
    if (price == 0) {
      // Reason 0 means unknown; this demo has no more specific standard reason.
      response.noBid(0);
    } else {
      var bid = BidResponseBuilder.bid("bid-1", imp.id(), price).banner().size(300, 250)
          .crid("creative-1").adomain("advertiser.example")
          .adm("<a href=\"https://advertiser.example/\"><img src=\"https://cdn.example/banner.png\" width=\"300\" height=\"250\"></a>")
          .build();
      response.addSeatBid("buyer-1", bid);
    }
    var res = response.build();
    var findings = BidCheck.response(requestView, res);
    if (!findings.ok()) throw new IllegalStateException(findings.errors().toString());
    // Demo policy: reject any warning, even when ok() is true.
    if (!findings.warnings().isEmpty()) {
      System.err.println("policy rejected bid: " + findings.warnings());
      res = BidResponseBuilder.create(requestView.auctionId()).currency("USD").noBid(0).build();
    }
    var payload = Json.toJsonBytes(res);
    var outputReport = Schema.response(payload);
    if (!outputReport.isOk()) throw new IllegalStateException(outputReport.toJson());
    return payload;
  }

  public static void main(String[] args) throws Exception {
    // SSP: construct an auction request and serialize it for the DSP.
    var req = BidRequestBuilder.create("auction-1").firstPrice().currency("USD").test()
        .site(Parts.site().id("site-1").domain("publisher.example").build())
        .addImp(ImpBuilders.banner("imp-1").size(300, 250).floor(1, "USD").secure().build())
        .build();
    var requestJson = Json.toJson(req);
    System.out.println("request: " + requestJson);
    for (double price : new double[] {2, 0.5, 0}) {
      var responseJson = new String(respond(requestJson, price), StandardCharsets.UTF_8);
      // SSP: decode the returned response and inspect its query view.
      var responseView = ResponseView.of(Json.parseBidResponse(responseJson));
      if (responseView.noBid() != (price < 1)) throw new IllegalStateException("unexpected decision");
      System.out.printf("price=%s no_bid=%s response=%s%n", price, responseView.noBid(), responseJson);
    }
  }
}
