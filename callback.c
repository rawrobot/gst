#include <stdio.h>
#include <stdlib.h>


#include "callback.h"
#include "goplugin/gstgocallback.h"


void X_gst_pad_set_chain_function(GstPad * pad) {
     gst_pad_set_chain_function(pad, go_callback_chain);
}

GstFlowReturn 
X_gst_pad_push(GstObject * parent, GstBuffer * buffer){
    GstGoCallback * filter = GST_GOCALLBACK (parent);
    return gst_pad_push( filter->srcpad, buffer) ;   
}

// @see https://developer.gnome.org/gobject/stable/gobject-properties.html
void x_go_object_set_property_guint64(GstElement *e, const char* prop, guint64 val) {
    GValue gval = G_VALUE_INIT;
    g_value_init (&gval, G_TYPE_UINT64);
    g_value_set_uint64 (&gval, val);
    printf("val %"G_GUINT64_FORMAT " ", val) ;
    g_object_set_property (G_OBJECT (e), prop, &gval);
    g_value_unset (&gval);
} 

static GstFlowReturn gst_govideocallback_transform_frame_ip (GstVideoFilter * filter,
    GstVideoFrame * frame){
    return go_transform_frame_ip(filter, frame) ;

}

void X_go_set_callback_transform_ip(GstElement *e) {
    guint64 val ;
    val = (guint64)gst_govideocallback_transform_frame_ip ;
    x_go_object_set_property_guint64(e, "transform-ip-callback", val);
} 

void X_go_set_callback_id(GstElement *e, guint64 val) {
    x_go_object_set_property_guint64(e, "caller-id", val);
} 