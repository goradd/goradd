// Adds an image capture function to manage an image capture control

goradd.imageCapture = function(controlId, selectName, shape, zoom, unsupportedError, data, mimeType, quality) {
    var $control = $j("#" + controlId);
    var timer;
    var $canvas = $control.find("canvas");
    var canvas = $canvas[0];
    var context = canvas.getContext('2d');
    var $button = $control.find("button");
    var btnText = $button.text();

    var $videoElement = $j("#imageCaptureVideo"); // reuse the element if its already here
    if ($videoElement.length == 0) {
        $videoElement = $j('<video playsinline></video>');
    }
    var videoElement = $videoElement[0];
    var videoDevice;

    if  (!('mediaDevices' in navigator)) {
        // either browser or device does not support media capture
        $control.find("button").disable();
        $control.attr("title", unsupportedError)
    }

    var base_image = new Image();
    base_image.src = data;
    base_image.onload = function(){
        canvas.getContext('2d').drawImage(base_image, 0, 0);
    };

    $control.on("imageCapture", function() {
        var d = $control.data("viewing");

        if (d) {
            clearInterval(timer);
            $control.data("viewing", false)
            $button.text(btnText);
            var newData = canvas.toDataURL("image/" + mimeType, quality).split(',')[1];
            goradd.setControlValue(controlId, "data", newData);
            videoDevice.stop();
        } else {
            // Beginning viewing the camera
            $control.data("viewing", true);
            $button.text(selectName);

            navigator.mediaDevices.getUserMedia({video: true}).then( function(mediaStream) {
                if (HTMLMediaElement) {
                    videoElement.srcObject = mediaStream;
                } else {
                    videoElement.src = URL.createObjectURL(mediaStream);
                }

                videoDevice = mediaStream.getVideoTracks()[0];

                var x,y, w, h;

                videoElement.onplay = function () {
                    timer = setInterval(function() {
                        try {
                            context.drawImage(videoElement, -x, -y, w, h);
                        } catch(err) {
                            console.debug(err)
                        }
                    }, 100);
                };

                videoElement.onloadedmetadata = function () {
                   // compute clip and resize values to center the image in the canvas
                    var hCanvas = canvas.height;
                    var wCanvas = canvas.width;
                    var hScale = videoElement.videoHeight / hCanvas;
                    var wScale = videoElement.videoWidth / wCanvas;
                    var scale = zoom / 100 + 1;
                    var smallestDimension;
                    if (wScale < hScale) {
                        w = scale * wCanvas;
                        h = scale * videoElement.videoHeight / wScale ;
                        x = 0;
                        y = scale * (h - hCanvas) / 2 ;
                    } else {
                        h = scale * hCanvas;
                        w = scale * videoElement.videoWidth / hScale;
                        y = 0;
                        x = scale * (w - wCanvas) / 2;

                    }
                    videoElement.play(); // Safari bug
                    var context = canvas.getContext('2d');

                    // flip image so selfie is easier to take
                    context.translate(canvas.width, 0);
                    context.scale(-1,1);

                    if (shape === "circle") {
                        context.beginPath();
                        context.rect(0, 0, wCanvas, hCanvas);
                        context.fillStyle = "white";
                        context.fill();

                        context.beginPath();
                        var minDim = Math.min(hCanvas, wCanvas);
                        context.arc(wCanvas / 2, hCanvas / 2, minDim / 2, 0, 2 * Math.PI, false);
                        context.clip();
                    }
                };


            }).catch(function(err) {
                $control.find("button").prop("disabled", true);
                $control.attr("title", unsupportedError)
                if (videoDevice) {
                    videoDevice.stop();
                }
            });
        }
    });
};

