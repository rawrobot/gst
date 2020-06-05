#include <stdio.h>
#include <stdlib.h>

#include "gst.h"

void X_gst_shim_init() {
  gchar *nano_str;
  guint major, minor, micro, nano;

  fprintf(stderr, "[ GSTREAMER ] shim init\n");
  gst_init(0, NULL);

  gst_version (&major, &minor, &micro, &nano);

  if (nano == 1)
    nano_str = "(CVS)";
  else if (nano == 2)
    nano_str = "(Prerelease)";
  else
    nano_str = "";

  printf ("[ GST ] program is linked against GStreamer %d.%d.%d %s\n",
          major, minor, micro, nano_str);

  return;
}

void X_gst_bin_add(GstElement *p, GstElement *element) {
  gst_bin_add(GST_BIN(p), element);

  return;
}

void X_gst_bin_remove(GstElement *p, GstElement *element) {
  gst_bin_remove(GST_BIN(p), element);

  return;
}

void X_gst_g_object_set_string(GstElement *e, const gchar* p_name, gchar* p_value) {
  g_object_set(G_OBJECT(e), p_name, p_value, NULL);
}

void X_gst_g_object_set_int(GstElement *e, const gchar* p_name, gint p_value) {
  g_object_set(G_OBJECT(e), p_name, p_value, NULL);
}

void X_gst_g_object_set_uint(GstElement *e, const gchar* p_name, guint p_value) {
  g_object_set(G_OBJECT(e), p_name, p_value, NULL);
}

void X_gst_g_object_set_uint64(GstElement *e, const gchar* p_name, guint64 p_value) {
  g_object_set(G_OBJECT(e), p_name, p_value, NULL);
}

void X_gst_g_object_set_bool(GstElement *e, const gchar* p_name, gboolean p_value) {
  g_object_set(G_OBJECT(e), p_name, p_value, NULL);
}

void X_gst_g_object_set_caps(GstElement *e, const gchar* p_name, const GstCaps *p_value) {
  g_object_set(G_OBJECT(e), p_name, p_value, NULL);
}

void X_gst_g_object_set_structure(GstElement *e, const gchar* p_name, const GstStructure *p_value) {
  g_object_set(G_OBJECT(e), p_name, p_value, NULL);
}

void X_gst_g_object_setv(GObject *object, guint n_properties, const gchar *names[], const GValue value[]) {
  //g_object_setv(object, n_properties, names, value);
}


void cb_new_pad(GstElement *element, GstPad *pad, gpointer data) {
  ElementUserData *d = (ElementUserData *)data;
  go_callback_new_pad(element, pad, d->callbackId);
}


void X_g_signal_connect(GstElement* element, gchar* detailed_signal, guint64 callbackId) {
  printf("[ GST ] g_signal_connect called with signal %s\n", detailed_signal);
  
  ElementUserData *d = calloc(1, sizeof(ElementUserData));
  d->callbackId = callbackId;

  g_signal_connect(element, detailed_signal, G_CALLBACK(cb_new_pad), d);
}

void X_g_signal_connect_data(gpointer instance, const gchar *detailed_signal, void (*f)(GstElement*, GstBus*, GstMessage*, gpointer), gpointer data, GClosureNotify destroy_data, GConnectFlags connect_flags) {
  printf("[ GST ] g_signal_connect_data called\n");
  g_signal_connect_data(instance, detailed_signal, G_CALLBACK(f), data, destroy_data, connect_flags);
}

GstElement *X_gst_bin_get_by_name(GstElement* element, const gchar* name) {
  GstElement *e = gst_bin_get_by_name(GST_BIN(element), name);
  if (e != NULL) {

  }
  return e;
}

GstElementClass *X_GST_ELEMENT_GET_CLASS(GstElement *element) {
  return GST_ELEMENT_GET_CLASS(element);
}

// Should set GST_DEBUG_DOT_DIR to output directory
// and run with --gst-enable-gst-debug command line switch
void X_GST_DEBUG_BIN_TO_DOT_FILE(GstElement *element, const gchar* name) {
  GST_DEBUG_BIN_TO_DOT_FILE(GST_BIN(element), GST_DEBUG_GRAPH_SHOW_ALL, name);
}

void X_g_signal_emit_buffer_by_name(GstElement* element, const gchar* detailed_signal, GstBuffer* buffer, GstFlowReturn* ret) {
  g_signal_emit_by_name(element, detailed_signal, buffer, ret);
}

GstBuffer *X_gst_buffer_new_wrapped(gchar* src, gsize len) {
  GstBuffer* dst;

  dst = gst_buffer_new_allocate(NULL, len, NULL);
  gst_buffer_fill(dst, 0, src, len);

  return dst;
}

gboolean X_gst_buffer_map(GstBuffer* gstBuffer, GstMapInfo* mapInfo) {
  return gst_buffer_map(gstBuffer, mapInfo, GST_MAP_READ);
}

void X_gst_pipeline_use_clock(GstElement *element, GstClock *clock) {
  gst_pipeline_use_clock(GST_PIPELINE(element), clock);
}

void X_gst_element_set_start_time_none(GstElement *element) {
  gst_element_set_start_time(element, GST_CLOCK_TIME_NONE);
}

void X_gst_structure_set_string(GstStructure *structure, const gchar *name, gchar *value) {
  GValue gv;
  memset(&gv, 0, sizeof(GValue));
  g_value_init(&gv, G_TYPE_STRING);
  g_value_set_string(&gv, value);
  gst_structure_set_value(structure, name, &gv);
}

void X_gst_structure_set_int(GstStructure *structure, const gchar *name, int value) {

  GValue gv;
  memset(&gv, 0, sizeof(GValue));
  g_value_init(&gv, G_TYPE_INT);
  g_value_set_int(&gv, value);
  gst_structure_set_value(structure, name, &gv);
}

void X_gst_structure_set_uint(GstStructure *structure, const gchar *name, guint value) {

  GValue gv;
  memset(&gv, 0, sizeof(GValue));
  g_value_init(&gv, G_TYPE_UINT);
  g_value_set_uint(&gv, value);
  gst_structure_set_value(structure, name, &gv);
}

void X_gst_structure_set_bool(GstStructure *structure, const gchar *name, gboolean value) {

  GValue gv;
  memset(&gv, 0, sizeof(GValue));
  g_value_init(&gv, G_TYPE_BOOLEAN);
  g_value_set_boolean(&gv, value);
  gst_structure_set_value(structure, name, &gv);
}

// events
GstEventType X_GST_EVENT_TYPE(GstEvent* event) {
    return GST_EVENT_CAST(event)->type;
}

// messages
GstMessageType X_GST_MESSAGE_TYPE(GstMessage *message) {
    return GST_MESSAGE_TYPE(message);
}

// bus
GstBus* X_gst_pipeline_get_bus(GstElement* element) {
	return gst_pipeline_get_bus(GST_PIPELINE(element));
}

GstClock * X_gst_pipeline_get_clock(GstElement* element) {
  return gst_pipeline_get_clock(GST_PIPELINE(element));
}


GstClockTime X_gst_pipeline_get_delay(GstElement* element) {
  return gst_pipeline_get_delay(GST_PIPELINE(element));
}


GstClockTime X_gst_pipeline_get_latency(GstElement* element) {
    return gst_pipeline_get_latency(GST_PIPELINE(element));
}

void X_gst_pipeline_set_latency(GstElement* element, GstClockTime clockTime) {
    gst_pipeline_set_latency(GST_PIPELINE(element), clockTime);
}


GstFlowReturn X_gst_app_src_push_buffer(GstElement* element, void *buffer,int len) {
    GstFlowReturn ret ; 
    gpointer p = g_memdup(buffer, len);
    GstBuffer *data = gst_buffer_new_wrapped(p, len); 
    ret  = gst_app_src_push_buffer(GST_APP_SRC(element), data);
    return ret ;
}

GstClockTime X_gst_buffer_get_duration(GstBuffer* buffer) {
  return GST_BUFFER_DURATION(buffer);
}

gchar* X_gst_pad_get_name(GstPad* pad) {
  return gst_pad_get_name(pad);
}

//By BKSWORM


void cb_bus_message(GstBus * bus, GstMessage * message, gpointer poll_data) {
  GError *err;
  gchar *debug_info;
 // BusMessageUserData *s = (BusMessageUserData *)user_data;
  //int pipelineId = s->pipelineId;

  /* Print error details on the screen */
  gst_message_parse_error (message, &err, &debug_info);
  g_printerr ("Error received from element %s: %s\n", GST_OBJECT_NAME (message->src), err->message);
  g_printerr ("Debugging information: %s\n", debug_info ? debug_info : "none");
  ///goHandleErrorMessage(msg, pipelineId, err, debug_info);
  g_clear_error (&err);
  g_free (debug_info);
  //go_callback_bus_message_thunk(bus, message, poll_data);
}

// Function to push go buf to pipe line buffer is duped
gboolean x_push_buffer_async(GstElement *element, void *buffer,int len, int framerate) {
    GstFlowReturn ret;
    gpointer p = g_memdup(buffer, len);
    GstBuffer *data = gst_buffer_new_wrapped(p, len); //TODO: do we need to free it?
    GST_BUFFER_DURATION (data) = gst_util_uint64_scale (1, GST_SECOND, framerate);


    // Push the buffer into the appsrc
    g_signal_emit_by_name (GST_APP_SRC(G_OBJECT(element)), "push-buffer", data, &ret);
   gst_buffer_unref (data) ;

    if (ret != GST_FLOW_OK) {
        // We got some error, stop sending data
        g_print ("push async error\n");
        return FALSE;
    }

    return TRUE;
}

void x_g_signal_emit_by_name(GstElement *element, const gchar *name){
	g_signal_emit_by_name(element, name, NULL);
}