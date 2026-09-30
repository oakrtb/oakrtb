package com.oakrtb.sdk.builder;

import com.oakrtb.openrtb.v2.Audio;
import com.oakrtb.openrtb.v2.Banner;
import com.oakrtb.openrtb.v2.Imp;
import com.oakrtb.openrtb.v2.Native;
import com.oakrtb.openrtb.v2.Video;

/**
 * Factories and fluent builders for each ad format's {@link Imp}.
 *
 * <p>Conveniently assembles banner, video, audio, and native Imps for {@link BidRequestBuilder#addImp(Imp)}.
 */
public final class ImpBuilders {
  private ImpBuilders() {}

  /**
   * Creates a Banner Imp builder.
   *
   * @param id unique Imp ID
   * @return a Banner Imp builder
   */
  public static BannerImp banner(String id) {
    return new BannerImp(id);
  }

  /**
   * Creates a Video Imp builder.
   *
   * @param id unique Imp ID
   * @return a Video Imp builder
   */
  public static VideoImp video(String id) {
    return new VideoImp(id);
  }

  /**
   * Creates an Audio Imp builder.
   *
   * @param id unique Imp ID
   * @return an Audio Imp builder
   */
  public static AudioImp audio(String id) {
    return new AudioImp(id);
  }

  /**
   * Creates a Native Imp builder.
   *
   * @param id unique Imp ID
   * @return a Native Imp builder
   */
  public static NativeImp nativeAd(String id) {
    return new NativeImp(id);
  }

  /**
   * Shared Imp builder base class for floor, security, tag, and other common fields.
   *
   * @param <T> concrete subtype (CRTP)
   */
  public abstract static class ImpBase<T extends ImpBase<T>> {
    protected final Imp.Builder imp = Imp.newBuilder();

    @SuppressWarnings("unchecked")
    protected T self() {
      return (T) this;
    }

    protected ImpBase(String id) {
      imp.setId(id);
    }

    /**
     * Sets the floor and currency.
     *
     * @param bidfloor price floor
     * @param cur ISO-4217 currency code
     * @return this builder
     */
    public T floor(double bidfloor, String cur) {
      imp.setBidfloor(bidfloor).setBidfloorcur(cur);
      return self();
    }

    /**
     * Requires HTTPS creatives (secure=1).
     *
     * @return this builder
     */
    public T secure() {
      imp.setSecure(1);
      return self();
    }

    /**
     * Sets the placement tag ID (tagid).
     *
     * @param tagid tag identifier
     * @return this builder
     */
    public T tagId(String tagid) {
      imp.setTagid(tagid);
      return self();
    }

    /**
     * Marks interstitial inventory (instl=1).
     *
     * @return this builder
     */
    public T interstitial() {
      imp.setInstl(1);
      return self();
    }

    /**
     * Marks rewarded inventory (rwdd=1).
     *
     * @return this builder
     */
    public T rewarded() {
      imp.setRwdd(1);
      return self();
    }

    /**
     * Builds a protobuf {@link Imp}.
     *
     * @return the impression object
     */
    public abstract Imp build();
  }

  /**
   * Banner Imp builder.
   */
  public static final class BannerImp extends ImpBase<BannerImp> {
    private final Banner.Builder banner = Banner.newBuilder();

    BannerImp(String id) {
      super(id);
    }

    /**
     * Sets the Banner width and height.
     *
     * @param w width in pixels
     * @param h height in pixels
     * @return this builder
     */
    public BannerImp size(int w, int h) {
      banner.setW(w).setH(h);
      return this;
    }

    /**
     * Sets the ad position (pos).
     *
     * @param pos OpenRTB ad position enum value
     * @return this builder
     */
    public BannerImp pos(int pos) {
      banner.setPos(pos);
      return this;
    }

    /**
     * Sets allowed MIME types.
     *
     * @param mimes MIME types, e.g. image/jpeg
     * @return this builder
     */
    public BannerImp mimes(String... mimes) {
      banner.clearMimes();
      for (String m : mimes) {
        banner.addMimes(m);
      }
      return this;
    }

    @Override
    public Imp build() {
      return imp.setBanner(banner).build();
    }
  }

  /**
   * Video Imp builder.
   */
  public static final class VideoImp extends ImpBase<VideoImp> {
    private final Video.Builder video = Video.newBuilder();

    VideoImp(String id) {
      super(id);
    }

    /**
     * Sets allowed MIME types (required).
     *
     * @param mimes MIME types, e.g. video/mp4
     * @return this builder
     */
    public VideoImp mimes(String... mimes) {
      video.clearMimes();
      for (String m : mimes) {
        video.addMimes(m);
      }
      return this;
    }

    /**
     * Sets the minimum and maximum duration in seconds.
     *
     * @param min minimum duration
     * @param max maximum duration
     * @return this builder
     */
    public VideoImp duration(int min, int max) {
      video.setMinduration(min).setMaxduration(max);
      return this;
    }

    /**
     * Sets supported VAST protocols.
     *
     * @param protocols protocol enum values
     * @return this builder
     */
    public VideoImp protocols(int... protocols) {
      video.clearProtocols();
      for (int p : protocols) {
        video.addProtocols(p);
      }
      return this;
    }

    /**
     * Sets player dimensions.
     *
     * @param w width
     * @param h height
     * @return this builder
     */
    public VideoImp size(int w, int h) {
      video.setW(w).setH(h);
      return this;
    }

    /**
     * Sets startdelay.
     *
     * @param v startdelay value
     * @return this builder
     */
    public VideoImp startDelay(int v) {
      video.setStartdelay(v);
      return this;
    }

    /**
     * Sets the video placement type (plcmt).
     *
     * @param plcmt plcmt enum value
     * @return this builder
     */
    public VideoImp plcmt(int plcmt) {
      video.setPlcmt(plcmt);
      return this;
    }

    /**
     * Sets linearity (linear/nonlinear).
     *
     * @param v linearity value
     * @return this builder
     */
    public VideoImp linearity(int v) {
      video.setLinearity(v);
      return this;
    }

    /**
     * Enables skipping and sets skipafter in seconds.
     *
     * @param skipafter seconds before skipping is allowed
     * @return this builder
     */
    public VideoImp skip(int skipafter) {
      video.setSkip(1).setSkipafter(skipafter);
      return this;
    }

    /**
     * Sets pod information.
     *
     * @param podid pod id
     * @param slotinpod slot within the pod
     * @return this builder
     */
    public VideoImp pod(String podid, int slotinpod) {
      video.setPodid(podid).setSlotinpod(slotinpod);
      return this;
    }

    /**
     * Sets playback methods (playbackmethod).
     *
     * @param methods playback method enum values
     * @return this builder
     */
    public VideoImp playbackMethod(int... methods) {
      video.clearPlaybackmethod();
      for (int m : methods) {
        video.addPlaybackmethod(m);
      }
      return this;
    }

    @Override
    public Imp build() {
      return imp.setVideo(video).build();
    }
  }

  /**
   * Audio Imp builder.
   */
  public static final class AudioImp extends ImpBase<AudioImp> {
    private final Audio.Builder audio = Audio.newBuilder();

    AudioImp(String id) {
      super(id);
    }

    /**
     * Sets allowed MIME types (required).
     *
     * @param mimes MIME types, e.g. audio/mpeg
     * @return this builder
     */
    public AudioImp mimes(String... mimes) {
      audio.clearMimes();
      for (String m : mimes) {
        audio.addMimes(m);
      }
      return this;
    }

    /**
     * Sets the minimum and maximum duration in seconds.
     *
     * @param min minimum duration
     * @param max maximum duration
     * @return this builder
     */
    public AudioImp duration(int min, int max) {
      audio.setMinduration(min).setMaxduration(max);
      return this;
    }

    /**
     * Sets supported protocols.
     *
     * @param protocols protocol enum values
     * @return this builder
     */
    public AudioImp protocols(int... protocols) {
      audio.clearProtocols();
      for (int p : protocols) {
        audio.addProtocols(p);
      }
      return this;
    }

    /**
     * Sets the feed type.
     *
     * @param feed feed enum value
     * @return this builder
     */
    public AudioImp feed(int feed) {
      audio.setFeed(feed);
      return this;
    }

    @Override
    public Imp build() {
      return imp.setAudio(audio).build();
    }
  }

  /**
   * Native Imp builder.
   */
  public static final class NativeImp extends ImpBase<NativeImp> {
    private final Native.Builder nativeAd = Native.newBuilder().setVer("1.2");

    NativeImp(String id) {
      super(id);
    }

    /**
     * Sets the Native Request JSON string (required).
     *
     * @param nativeRequestJson Native 1.x request object JSON
     * @return this builder
     */
    public NativeImp request(String nativeRequestJson) {
      nativeAd.setRequest(nativeRequestJson);
      return this;
    }

    /**
     * Sets the Native specification version (ver).
     *
     * @param ver version number, default 1.2
     * @return this builder
     */
    public NativeImp ver(String ver) {
      nativeAd.setVer(ver);
      return this;
    }

    @Override
    public Imp build() {
      return imp.setNative(nativeAd).build();
    }
  }
}
