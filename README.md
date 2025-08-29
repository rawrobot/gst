# gst

It is a small and simple go gstreamer binding. It is a fork from https://github.com/notedit/gst
Moreover it implements interface to gstreamer videofilter by intercepting xx_transform_ip() call. 


## NOTE

It is my very old project. It is only for portfolio purposes, to show children of Agile that a software architect can work with different programming languages.



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

### plugin/line
It is simple gstreamer "videofilter" example. It generates white box jpeg image with black line at the midle.

### plugin/movingline
This one shows video moving from up to down black line. 



