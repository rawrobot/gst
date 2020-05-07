#include <stdio.h>
#include <stdlib.h>


#include "callback.h"
#include "goplugin/gstgocallback.h"


void X_gst_pad_set_chain_function(GstPad * pad) {
    gst_pad_set_chain_function(pad, go_callback_chain);
}

GstFlowReturn
X_gst_pad_push(GstObject * parent, GstBuffer * buffer) {
    GstElement * filter = GST_ELEMENT_CAST (parent);
    GList* srcpads = filter->srcpads ;
    GstPad* srcpad = GST_PAD_CAST(g_list_nth_data(srcpads, 0)) ;
    return gst_pad_push( srcpad, buffer) ;
}

// @see https://developer.gnome.org/gobject/stable/gobject-properties.html
void x_go_object_set_property_guint64(GstElement *e, const char* prop, guint64 val) {
    GValue gval = G_VALUE_INIT;
    g_value_init (&gval, G_TYPE_UINT64);
    g_value_set_uint64 (&gval, val);
    g_object_set_property (G_OBJECT (e), prop, &gval);
    g_value_unset (&gval);
}

void x_go_object_set_property_gpointer(GstElement *e, const char* prop, void* val) {
    GValue gval = G_VALUE_INIT;
    g_value_init (&gval, G_TYPE_POINTER);
    g_value_set_pointer (&gval, val);
    g_object_set_property (G_OBJECT (e), prop, &gval);
    g_value_unset (&gval);
}


static GstFlowReturn gst_govideocallback_transform_frame_ip (GstVideoFilter * filter,
        GstVideoFrame * frame) {
    return go_transform_frame_ip(filter, frame) ;

}

void X_go_set_callback_transform_ip(GstElement *e) {
    void*  val ;
    val = (void*)gst_govideocallback_transform_frame_ip ;
    x_go_object_set_property_gpointer(e, "transform-ip-callback", val);
}

void X_go_set_callback_id(GstElement *e, guint64 val) {
    x_go_object_set_property_guint64(e, "caller-id", val);
}

guint8* get_frame_data(GstVideoFrame* frame, int plane ) {
    return GST_VIDEO_FRAME_PLANE_DATA(frame, plane) ;
}
guint32 get_frame_stride (GstVideoFrame* frame, int plane ) {
    return GST_VIDEO_FRAME_PLANE_STRIDE(frame, plane) ;
}
guint32 get_frame_pixel_stride(GstVideoFrame* frame, int plane ) {
    return GST_VIDEO_FRAME_COMP_PSTRIDE(frame, plane) ;
}
guint32 get_frame_h(GstVideoFrame* frame, int plane ) {
    return GST_VIDEO_FRAME_COMP_HEIGHT(frame, plane) ;
}

guint32 get_frame_w(GstVideoFrame* frame, int plane ) {
    return GST_VIDEO_FRAME_COMP_WIDTH(frame, plane) ;
}

guint32 get_frame_format(GstVideoFrame* frame) {
    return GST_VIDEO_FRAME_FORMAT(frame) ;
}

guint32 get_frame_n_planes(GstVideoFrame* frame) {
    return  GST_VIDEO_FRAME_N_PLANES (frame) ;
} 