#ifndef __GO__GST_CALLBACK_H__
#define __GO__GST_CALLBACK_H__

#include <stdlib.h>
#include <gst/gst.h>
#include <string.h>
#include <gst/base/gstbasetransform.h>
#include <gst/video/gstvideofilter.h>

extern GstFlowReturn go_callback_chain(GstPad * pad, GstObject * parent, GstBuffer * buffer) ;

extern void X_gst_pad_set_chain_function(GstPad * pad) ;
extern GstFlowReturn X_gst_pad_push(GstObject * parent, GstBuffer * buffer) ;

extern GstFlowReturn go_transform_frame_ip (GstVideoFilter * filter,
   GstVideoFrame * frame) ;
extern void X_go_set_callback_transform_ip(GstElement *e) ;
extern void X_go_set_callback_id(GstElement *e, guint64 val)  ;

extern guint8* get_frame_data(GstVideoFrame* frame, int plane );
extern guint32 get_frame_stride (GstVideoFrame* frame, int plane );
extern guint32 get_frame_pixel_stride(GstVideoFrame* frame, int plane );
extern guint32 get_frame_h(GstVideoFrame* frame, int plane );
extern guint32 get_frame_w(GstVideoFrame* frame, int plane );
extern guint32 get_frame_format(GstVideoFrame* frame) ;
extern guint32 get_frame_n_planes(GstVideoFrame* frame) ;

#endif