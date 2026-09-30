package com.oakrtb.sdk.builder;

import com.oakrtb.openrtb.v2.App;
import com.oakrtb.openrtb.v2.Content;
import com.oakrtb.openrtb.v2.Device;
import com.oakrtb.openrtb.v2.Dooh;
import com.oakrtb.openrtb.v2.Geo;
import com.oakrtb.openrtb.v2.Publisher;
import com.oakrtb.openrtb.v2.Site;

/**
 * Lightweight builders for OpenRTB subobjects such as Site, App, Dooh, and Device.
 *
 * <p>Works with {@link BidRequestBuilder} to assemble inventory, device, publisher, and other fields.
 */
public final class Parts {
  private Parts() {}

  /** Creates a Site builder. */
  public static SiteBuilder site() {
    return new SiteBuilder();
  }

  /** Creates an App builder. */
  public static AppBuilder app() {
    return new AppBuilder();
  }

  /** Creates a Dooh builder. */
  public static DoohBuilder dooh() {
    return new DoohBuilder();
  }

  /** Creates a Device builder. */
  public static DeviceBuilder device() {
    return new DeviceBuilder();
  }

  /** Creates a Publisher builder. */
  public static PublisherBuilder publisher() {
    return new PublisherBuilder();
  }

  /** Creates a Content builder. */
  public static ContentBuilder content() {
    return new ContentBuilder();
  }

  /** Creates a Geo builder. */
  public static GeoBuilder geo() {
    return new GeoBuilder();
  }

  /** Fluent builder for a website (Site). */
  public static final class SiteBuilder {
    private final Site.Builder b = Site.newBuilder();

    /**
     * Sets the site ID.
     *
     * @param id site identifier
     * @return this builder
     */
    public SiteBuilder id(String id) {
      b.setId(id);
      return this;
    }

    /**
     * Sets the site name.
     *
     * @param name name
     * @return this builder
     */
    public SiteBuilder name(String name) {
      b.setName(name);
      return this;
    }

    /**
     * Sets the site domain.
     *
     * @param domain domain
     * @return this builder
     */
    public SiteBuilder domain(String domain) {
      b.setDomain(domain);
      return this;
    }

    /**
     * Sets the current page URL.
     *
     * @param page page URL
     * @return this builder
     */
    public SiteBuilder page(String page) {
      b.setPage(page);
      return this;
    }

    /**
     * Sets content categories (cat), such as IAB categories.
     *
     * @param cats category list
     * @return this builder
     */
    public SiteBuilder cat(String... cats) {
      b.clearCat();
      for (String c : cats) {
        b.addCat(c);
      }
      return this;
    }

    /**
     * Sets the publisher (publisher).
     *
     * @param p Publisher object
     * @return this builder
     */
    public SiteBuilder publisher(Publisher p) {
      b.setPublisher(p);
      return this;
    }

    /**
     * Builds a {@link Site}.
     *
     * @return a Site object
     */
    public Site build() {
      return b.build();
    }
  }

  /** Fluent builder for a mobile app (App). */
  public static final class AppBuilder {
    private final App.Builder b = App.newBuilder();

    /**
     * Sets the app ID.
     *
     * @param id app identifier
     * @return this builder
     */
    public AppBuilder id(String id) {
      b.setId(id);
      return this;
    }

    /**
     * Sets the app name.
     *
     * @param name name
     * @return this builder
     */
    public AppBuilder name(String name) {
      b.setName(name);
      return this;
    }

    /**
     * Sets the app bundle ID (bundle).
     *
     * @param bundle bundle ID
     * @return this builder
     */
    public AppBuilder bundle(String bundle) {
      b.setBundle(bundle);
      return this;
    }

    /**
     * Sets the app domain.
     *
     * @param domain domain
     * @return this builder
     */
    public AppBuilder domain(String domain) {
      b.setDomain(domain);
      return this;
    }

    /**
     * Sets the publisher (publisher).
     *
     * @param p Publisher object
     * @return this builder
     */
    public AppBuilder publisher(Publisher p) {
      b.setPublisher(p);
      return this;
    }

    /**
     * Sets content context (content).
     *
     * @param c Content object
     * @return this builder
     */
    public AppBuilder content(Content c) {
      b.setContent(c);
      return this;
    }

    /**
     * Builds an {@link App}.
     *
     * @return an App object
     */
    public App build() {
      return b.build();
    }
  }

  /** Fluent builder for digital out-of-home inventory (Dooh). */
  public static final class DoohBuilder {
    private final Dooh.Builder b = Dooh.newBuilder();

    /**
     * Sets the DOOH placement ID.
     *
     * @param id identifier
     * @return this builder
     */
    public DoohBuilder id(String id) {
      b.setId(id);
      return this;
    }

    /**
     * Sets the DOOH placement name.
     *
     * @param name name
     * @return this builder
     */
    public DoohBuilder name(String name) {
      b.setName(name);
      return this;
    }

    /**
     * Sets venue types (venuetype).
     *
     * @param ids venue type ID list
     * @return this builder
     */
    public DoohBuilder venueType(String... ids) {
      b.clearVenuetype();
      for (String id : ids) {
        b.addVenuetype(id);
      }
      return this;
    }

    /**
     * Sets the venue type taxonomy (venuetypetax).
     *
     * @param tax taxonomy value
     * @return this builder
     */
    public DoohBuilder venueTypeTax(int tax) {
      b.setVenuetypetax(tax);
      return this;
    }

    /**
     * Sets the publisher (publisher).
     *
     * @param p Publisher object
     * @return this builder
     */
    public DoohBuilder publisher(Publisher p) {
      b.setPublisher(p);
      return this;
    }

    /**
     * Builds a {@link Dooh}.
     *
     * @return a Dooh object
     */
    public Dooh build() {
      return b.build();
    }
  }

  /** Fluent builder for a device (Device). */
  public static final class DeviceBuilder {
    private final Device.Builder b = Device.newBuilder();

    /**
     * Sets the User-Agent.
     *
     * @param ua UA string
     * @return this builder
     */
    public DeviceBuilder ua(String ua) {
      b.setUa(ua);
      return this;
    }

    /**
     * Sets the IPv4 address.
     *
     * @param ip IP address
     * @return this builder
     */
    public DeviceBuilder ip(String ip) {
      b.setIp(ip);
      return this;
    }

    /**
     * Sets the device type (devicetype).
     *
     * @param t OpenRTB device type enum value
     * @return this builder
     */
    public DeviceBuilder deviceType(int t) {
      b.setDevicetype(t);
      return this;
    }

    /**
     * Sets the device manufacturer (make).
     *
     * @param make manufacturer
     * @return this builder
     */
    public DeviceBuilder make(String make) {
      b.setMake(make);
      return this;
    }

    /**
     * Sets the device model (model).
     *
     * @param model model
     * @return this builder
     */
    public DeviceBuilder model(String model) {
      b.setModel(model);
      return this;
    }

    /**
     * Sets the operating system and version.
     *
     * @param os operating system name
     * @param osv operating system version
     * @return this builder
     */
    public DeviceBuilder os(String os, String osv) {
      b.setOs(os).setOsv(osv);
      return this;
    }

    /**
     * Sets the IFA (advertising identifier).
     *
     * @param ifa IFA string
     * @return this builder
     */
    public DeviceBuilder ifa(String ifa) {
      b.setIfa(ifa);
      return this;
    }

    /**
     * Sets geolocation (geo).
     *
     * @param g Geo object
     * @return this builder
     */
    public DeviceBuilder geo(Geo g) {
      b.setGeo(g);
      return this;
    }

    /**
     * Builds a {@link Device}.
     *
     * @return a Device object
     */
    public Device build() {
      return b.build();
    }
  }

  /** Fluent builder for a publisher (Publisher). */
  public static final class PublisherBuilder {
    private final Publisher.Builder b = Publisher.newBuilder();

    /**
     * Sets the publisher ID.
     *
     * @param id identifier
     * @return this builder
     */
    public PublisherBuilder id(String id) {
      b.setId(id);
      return this;
    }

    /**
     * Sets the publisher name.
     *
     * @param name name
     * @return this builder
     */
    public PublisherBuilder name(String name) {
      b.setName(name);
      return this;
    }

    /**
     * Sets the publisher domain.
     *
     * @param domain domain
     * @return this builder
     */
    public PublisherBuilder domain(String domain) {
      b.setDomain(domain);
      return this;
    }

    /**
     * Builds a {@link Publisher}.
     *
     * @return a Publisher object
     */
    public Publisher build() {
      return b.build();
    }
  }

  /** Fluent builder for content (Content). */
  public static final class ContentBuilder {
    private final Content.Builder b = Content.newBuilder();

    /**
     * Sets the content title.
     *
     * @param title title
     * @return this builder
     */
    public ContentBuilder title(String title) {
      b.setTitle(title);
      return this;
    }

    /**
     * Sets the series name (series).
     *
     * @param series series
     * @return this builder
     */
    public ContentBuilder series(String series) {
      b.setSeries(series);
      return this;
    }

    /**
     * Sets the season (season).
     *
     * @param season season identifier
     * @return this builder
     */
    public ContentBuilder season(String season) {
      b.setSeason(season);
      return this;
    }

    /**
     * Sets the episode number (episode).
     *
     * @param n episode number
     * @return this builder
     */
    public ContentBuilder episode(int n) {
      b.setEpisode(n);
      return this;
    }

    /**
     * Sets the content context (context).
     *
     * @param v context enum value
     * @return this builder
     */
    public ContentBuilder context(int v) {
      b.setContext(v);
      return this;
    }

    /**
     * Sets whether the content is live (livestream).
     *
     * @param v 0/1
     * @return this builder
     */
    public ContentBuilder livestream(int v) {
      b.setLivestream(v);
      return this;
    }

    /**
     * Sets whether the content is real-time (realtime).
     *
     * @param v 0/1
     * @return this builder
     */
    public ContentBuilder realtime(int v) {
      b.setRealtime(v);
      return this;
    }

    /**
     * Builds a {@link Content}.
     *
     * @return a Content object
     */
    public Content build() {
      return b.build();
    }
  }

  /** Fluent builder for geolocation (Geo). */
  public static final class GeoBuilder {
    private final Geo.Builder b = Geo.newBuilder();

    /**
     * Sets latitude and longitude.
     *
     * @param lat latitude
     * @param lon longitude
     * @return this builder
     */
    public GeoBuilder latLon(double lat, double lon) {
      b.setLat(lat).setLon(lon);
      return this;
    }

    /**
     * Sets the location type (type).
     *
     * @param t type enum value
     * @return this builder
     */
    public GeoBuilder type(int t) {
      b.setType(t);
      return this;
    }

    /**
     * Sets the country code (ISO-3166-1 alpha-3).
     *
     * @param c country code
     * @return this builder
     */
    public GeoBuilder country(String c) {
      b.setCountry(c);
      return this;
    }

    /**
     * Sets the city name.
     *
     * @param c city
     * @return this builder
     */
    public GeoBuilder city(String c) {
      b.setCity(c);
      return this;
    }

    /**
     * Builds a {@link Geo}.
     *
     * @return a Geo object
     */
    public Geo build() {
      return b.build();
    }
  }
}
