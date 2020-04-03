# gst

It is a small and simple go gstreamer binding. It is a fork from https://github.com/notedit/gst

## Install

*Ubuntu or Dedian*

```shell
apt-get install pkg-config
apt-get install libgstreamer1.0-dev libgstreamer-plugins-base1.0-dev libgstreamer-plugins-good1.0-dev gstreamer1.0-libav
```


## Examples

The are in the ./examples folder. Please be awere that exmples require gtk3 library and gtksink gstreamer plugins to be installed.
```shell
apt-get install libgtk-3-dev
apt-get install gstreamer1.0-gtk3
```

### gtksink

It shows video output from /dev/video0 at the gtk window.

### camui

It does the same thing but may pause/resume the video stream and take snapshots from the camera.


