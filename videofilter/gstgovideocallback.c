/* GStreamer
 * Copyright (C) 2020 FIXME <fixme@example.com>
 *
 * This library is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Library General Public
 * License as published by the Free Software Foundation; either
 * version 2 of the License, or (at your option) any later version.
 *
 * This library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 * Library General Public License for more details.
 *
 * You should have received a copy of the GNU Library General Public
 * License along with this library; if not, write to the
 * Free Software Foundation, Inc., 51 Franklin Street, Suite 500,
 * Boston, MA 02110-1335, USA.
 */
/**
 * SECTION:element-gstgovideocallback
 *
 * The govideocallback element does FIXME stuff.
 *
 * <refsect2>
 * <title>Example launch line</title>
 * |[
 * gst-launch-1.0 -v fakesrc ! govideocallback ! FIXME ! fakesink
 * ]|
 * FIXME Describe what the pipeline does.
 * </refsect2>
 */

#ifdef HAVE_CONFIG_H
#include "config.h"
#endif

#include <gst/gst.h>
#include <gst/video/video.h>
#include <gst/base/gstbasetransform.h>
#include <gst/video/gstvideofilter.h>
#include "gstgovideocallback.h"

GST_DEBUG_CATEGORY_STATIC (gst_govideocallback_debug_category);
#define GST_CAT_DEFAULT gst_govideocallback_debug_category

#define DEFAULT_HORIZONTAL_SPEED 10

/* prototypes */


static void gst_govideocallback_set_property (GObject * object,
    guint property_id, const GValue * value, GParamSpec * pspec);
static void gst_govideocallback_get_property (GObject * object,
    guint property_id, GValue * value, GParamSpec * pspec);
static void gst_govideocallback_dispose (GObject * object);
static void gst_govideocallback_finalize (GObject * object);

static gboolean gst_govideocallback_start (GstBaseTransform * trans);
static gboolean gst_govideocallback_stop (GstBaseTransform * trans);
static gboolean gst_govideocallback_set_info (GstVideoFilter * filter, GstCaps * incaps,
    GstVideoInfo * in_info, GstCaps * outcaps, GstVideoInfo * out_info);
static GstFlowReturn gst_govideocallback_transform_frame (GstVideoFilter * filter,
    GstVideoFrame * inframe, GstVideoFrame * outframe);
static GstFlowReturn gst_govideocallback_transform_frame_ip (GstVideoFilter * filter,
    GstVideoFrame * frame);

enum
{
  PROP_0,
  PROP_IP_CALLBACK,
  PROP_CALLER_ID
};

/* pad templates */

/* FIXME: add/remove formats you can handle */
#define VIDEO_SRC_CAPS \
    GST_VIDEO_CAPS_MAKE("{ I420, Y444, Y42B, UYVY, RGBA }")

/* FIXME: add/remove formats you can handle */
#define VIDEO_SINK_CAPS \
    GST_VIDEO_CAPS_MAKE("{ I420, Y444, Y42B, UYVY, RGBA }")


/* class initialization */

G_DEFINE_TYPE_WITH_CODE (GstGoVideoCallback, gst_govideocallback, GST_TYPE_VIDEO_FILTER,
  GST_DEBUG_CATEGORY_INIT (gst_govideocallback_debug_category, "govideocallback", 0,
  "debug category for govideocallback element"));

static void
gst_govideocallback_class_init (GstGoVideoCallbackClass * klass)
{
  GObjectClass *gobject_class = G_OBJECT_CLASS (klass);
  GstBaseTransformClass *base_transform_class = GST_BASE_TRANSFORM_CLASS (klass);
  GstVideoFilterClass *video_filter_class = GST_VIDEO_FILTER_CLASS (klass);

  /* Setting up pads and setting metadata should be moved to
     base_class_init if you intend to subclass this class. */
  gst_element_class_add_pad_template (GST_ELEMENT_CLASS(klass),
      gst_pad_template_new ("src", GST_PAD_SRC, GST_PAD_ALWAYS,
        gst_caps_from_string (VIDEO_SRC_CAPS)));
  gst_element_class_add_pad_template (GST_ELEMENT_CLASS(klass),
      gst_pad_template_new ("sink", GST_PAD_SINK, GST_PAD_ALWAYS,
        gst_caps_from_string (VIDEO_SINK_CAPS)));

  gst_element_class_set_static_metadata (GST_ELEMENT_CLASS(klass),
      "FIXME Long name", "Generic", "FIXME Description",
      "FIXME <fixme@example.com>");

  gobject_class->set_property = gst_govideocallback_set_property;
  gobject_class->get_property = gst_govideocallback_get_property;
  gobject_class->dispose = gst_govideocallback_dispose;
  gobject_class->finalize = gst_govideocallback_finalize;
  base_transform_class->start = GST_DEBUG_FUNCPTR (gst_govideocallback_start);
  base_transform_class->stop = GST_DEBUG_FUNCPTR (gst_govideocallback_stop);
  video_filter_class->set_info = GST_DEBUG_FUNCPTR (gst_govideocallback_set_info);
  video_filter_class->transform_frame = NULL ; //GST_DEBUG_FUNCPTR (gst_govideocallback_transform_frame);
  video_filter_class->transform_frame_ip = GST_DEBUG_FUNCPTR (gst_govideocallback_transform_frame_ip);

  g_object_class_install_property (gobject_class, PROP_IP_CALLBACK,
      g_param_spec_pointer ("transform-ip-callback", "Ip trasform callback",
          "Go lang callback for ip transformation",
          G_PARAM_READWRITE)) ; // |  G_PARAM_STATIC_STRINGS));  

g_object_class_install_property (gobject_class, PROP_CALLER_ID,
      g_param_spec_uint64 ("caller-id", "Ip trasform caller id",
          "It's ID golang callback provider",
           0, G_MAXUINT64, 0, 
          G_PARAM_READWRITE )) ;//| G_PARAM_STATIC_STRINGS));  

}

static void
gst_govideocallback_init (GstGoVideoCallback *govideocallback)
{
   
}

void
gst_govideocallback_set_property (GObject * object, guint property_id,
    const GValue * value, GParamSpec * pspec)
{
    guint64 v ;
  GstGoVideoCallback *govideocallback = GST_GOVIDEOCALLBACK (object);

  GST_DEBUG_OBJECT (govideocallback, "set_property");

  switch (property_id) {
      case PROP_IP_CALLBACK:
        v = g_value_get_uint64 (value) ;
        GST_DEBUG_OBJECT (govideocallback, "set_property transform-ip-callback %"G_GUINT64_FORMAT " ",v);
        govideocallback->ip_callback =  (gst_govideocallback_tarnsform_ip_t)v;
      break; 
      case PROP_CALLER_ID:
    v = g_value_get_uint64 (value); 
       GST_DEBUG_OBJECT (govideocallback, "set_property caller-id %"G_GUINT64_FORMAT " ",v);
      govideocallback->caller_id = v ;
      break; 
    default:
      G_OBJECT_WARN_INVALID_PROPERTY_ID (object, property_id, pspec);
      break;
  }
}

void
gst_govideocallback_get_property (GObject * object, guint property_id,
    GValue * value, GParamSpec * pspec)
{
  GstGoVideoCallback *govideocallback = GST_GOVIDEOCALLBACK (object);

  GST_DEBUG_OBJECT (govideocallback, "get_property");

  switch (property_id) {
      case PROP_IP_CALLBACK:
         g_value_set_uint64 (value, (guint64)govideocallback->ip_callback);
      break; 
    case PROP_CALLER_ID:
      g_value_set_uint64 (value, govideocallback->caller_id);
      break;  
    default:
      G_OBJECT_WARN_INVALID_PROPERTY_ID (object, property_id, pspec);
      break;
  }
}

void
gst_govideocallback_dispose (GObject * object)
{
  GstGoVideoCallback *govideocallback = GST_GOVIDEOCALLBACK (object);

  GST_DEBUG_OBJECT (govideocallback, "dispose");

  /* clean up as possible.  may be called multiple times */

  G_OBJECT_CLASS (gst_govideocallback_parent_class)->dispose (object);
}

void
gst_govideocallback_finalize (GObject * object)
{
  GstGoVideoCallback *govideocallback = GST_GOVIDEOCALLBACK (object);

  GST_DEBUG_OBJECT (govideocallback, "finalize");

  /* clean up object here */

  G_OBJECT_CLASS (gst_govideocallback_parent_class)->finalize (object);
}

static gboolean
gst_govideocallback_start (GstBaseTransform * trans)
{
  GstGoVideoCallback *govideocallback = GST_GOVIDEOCALLBACK (trans);

  GST_DEBUG_OBJECT (govideocallback, "start");

  return TRUE;
}

static gboolean
gst_govideocallback_stop (GstBaseTransform * trans)
{
  GstGoVideoCallback *govideocallback = GST_GOVIDEOCALLBACK (trans);

  GST_DEBUG_OBJECT (govideocallback, "stop");

  return TRUE;
}

static gboolean
gst_govideocallback_set_info (GstVideoFilter * filter, GstCaps * incaps,
    GstVideoInfo * in_info, GstCaps * outcaps, GstVideoInfo * out_info)
{
  GstGoVideoCallback *govideocallback = GST_GOVIDEOCALLBACK (filter);

  GST_DEBUG_OBJECT (govideocallback, "set_info");

  return TRUE;
}

/* transform */
static GstFlowReturn
gst_govideocallback_transform_frame (GstVideoFilter * filter, GstVideoFrame * inframe,
    GstVideoFrame * outframe)
{
  GstGoVideoCallback *govideocallback = GST_GOVIDEOCALLBACK (filter);

  GST_DEBUG_OBJECT (govideocallback, "transform_frame");

  return GST_FLOW_OK;
}

static GstFlowReturn
gst_govideocallback_transform_frame_ip (GstVideoFilter * filter, GstVideoFrame * frame)
{
    GstGoVideoCallback *govideocallback = GST_GOVIDEOCALLBACK (filter);
    GST_DEBUG_OBJECT (govideocallback, "transform_frame_ip");
    GST_DEBUG_OBJECT(govideocallback, "caller-id %"G_GUINT64_FORMAT" ", (govideocallback->caller_id));
  /*
  GstVideoInfo * video_info= NULL;
  GstBuffer * video_buffer= NULL;
   // set RGB pixels to black one at a time
   if (gst_video_frame_map (frame, video_info, video_buffer, GST_MAP_WRITE)) {
     guint8 *pixels = GST_VIDEO_FRAME_PLANE_DATA (frame, 0);
     guint stride = GST_VIDEO_FRAME_PLANE_STRIDE (frame, 0);
     guint pixel_stride = GST_VIDEO_FRAME_COMP_PSTRIDE (frame, 0);
    guint height =GST_VIDEO_FRAME_COMP_HEIGHT (frame, 0);
    guint width = GST_VIDEO_FRAME_COMP_WIDTH (frame, 0);
     guint  h, w ;
     for (h = 0; h < height; ++h) {
       for (w = 0; w < width; ++w) {
         guint8 *pixel = pixels + h * stride + w * pixel_stride;

         memset (pixel, 0, pixel_stride);
       }
     }

     gst_video_frame_unmap (frame);
   }

*/  
  if (govideocallback->ip_callback != NULL){
    GST_DEBUG_OBJECT (govideocallback, "transform_frame_ip call callback");
       govideocallback->ip_callback(filter, frame) ;
    
    } 
 
  return GST_FLOW_OK;
}

static gboolean
plugin_init (GstPlugin * plugin)
{

  /* FIXME Remember to set the rank if it's an element that is meant
     to be autoplugged by decodebin. */
  return gst_element_register (plugin, "govideocallback", GST_RANK_NONE,
      GST_TYPE_GOVIDEOCALLBACK);
}

/* FIXME: these are normally defined by the GStreamer build system.
   If you are creating an element to be included in gst-plugins-*,
   remove these, as they're always defined.  Otherwise, edit as
   appropriate for your external plugin package. */
#ifndef VERSION
#define VERSION "0.0.FIXME"
#endif
#ifndef PACKAGE
#define PACKAGE "FIXME_package"
#endif
#ifndef PACKAGE_NAME
#define PACKAGE_NAME "FIXME_package_name"
#endif
#ifndef GST_PACKAGE_ORIGIN
#define GST_PACKAGE_ORIGIN "http://FIXME.org/"
#endif

GST_PLUGIN_DEFINE (GST_VERSION_MAJOR,
    GST_VERSION_MINOR,
    govideocallback,
    "FIXME plugin description",
    plugin_init, VERSION, "LGPL", PACKAGE_NAME, GST_PACKAGE_ORIGIN)

