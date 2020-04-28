project('gstgovideocallback', 'c', version : '1.16.2', license : 'LGPL')
pkg = import('pkgconfig')

prj_name = meson.project_name()

plugins_install_dir = join_paths(get_option('libdir'), 'gstreamer-1.0')

cc = meson.get_compiler('c')

gst_version = meson.project_version()

api_version = '1.0'

gst_dep = dependency('gstreamer-1.0',
    fallback : ['gstreamer', 'gst_dep'])

gst_video_dep = dependency('gstreamer-video-1.0')
plugin_c_args = ['-DHAVE_CONFIG_H']

cdata = configuration_data()
cdata.set_quoted('PACKAGE_VERSION', gst_version)
cdata.set_quoted('PACKAGE', prj_name + '-plugin')
cdata.set_quoted('GST_LICENSE', meson.project_license())
cdata.set_quoted('GST_API_VERSION', api_version)
cdata.set_quoted('GST_PACKAGE_NAME', 'GStreamer gogst Plug-ins')
cdata.set_quoted('GST_PACKAGE_ORIGIN', 'https://github.com/bksworm/gst')
configure_file(output : 'config.h', configuration : cdata)

plugin_sources = [
  'gstgovideocallback.c'
  ]    
                                                               
gstplugin = library(prj_name,
  plugin_sources,
  c_args: plugin_c_args,
  dependencies : [gst_dep, gst_video_dep],
  install : true,
  install_dir : plugins_install_dir,
)

pkg.generate(gstplugin)
