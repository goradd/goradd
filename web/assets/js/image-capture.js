// Widget to manage the imageCapture control.
jQuery(function( $, undefined ) {

$.widget( "goradd.imageCapture", {
    options: {
        selectButtonName:    "Capture",          // Name that the button gets when in capture mode
        shape:               "rect",             // Mask shape
        zoom:                0,                  // % of zoom. 100 means 2x zoom. 200 means 3x zoom.
        data:                "data:image/jpeg;base64,",   // Mime typed and base64'd data.
        mimeType:            "jpeg",            // desired mime type
        quality:             0.92               // Quality setting for jpeg and webP
    },
    _mimeType: "image/jpeg",
    _viewing: false,
    _create: function() {
        var $control = this.element,
            id = $control.attr('id');

        this._super();

        if  (!('mediaDevices' in navigator)) {
            // either browser or device does not support media capture
            this._showErr();
            return;
        }

        this._videoElement = $j('<video playsinline></video>')[0];
        this._canvas = $control.find("canvas")[0];
        this._context = this._canvas.getContext('2d');
        this._captureButton = $control.find("#" + id + "_capture");
        this._switchButton = $control.find("#" + id + "_switch");

        this._originalButtonText = this._captureButton.text();
        this._on( this._captureButton, {
            "click": this._imageButtonClick
        });
        this._on( this._switchButton, {
            "click": this._switchButtonClick
        });

        this._on( this._videoElement, {
            "playing": this._playing,
            "loadedmetadata": this._loadmetadata
        });
        this._loadNewData();

        this._devices = null;
        if (navigator.mediaDevices.enumerateDevices) {
            navigator.mediaDevices.enumerateDevices()
                .then(function(devices) {
                    self._devices = devices;
                })
                .catch(function(err) {
                    self._showErr(err)
                });
        }
        this._deviceNum = 0;
    },
    _showErr: function(err) {
        var $control = this.element,
            id = $control.attr('id');

        this._captureButton.prop("disabled", true);
        var p = $control.find("#" + id + "_err");
        p.show();
        if (err) {
            p.text(err);
        }
    },
    _loadNewData: function() {
        var self = this;
        var base_image = new Image();
        base_image.src = this.options.data;
        base_image.onload = function(){
            self._context.drawImage(base_image, 0, 0);
        };
    },
    _imageButtonClick: function() {
        if (this._viewing) {
            this._capture();
        } else {
            this._view();
        }
    },
    turnOff: function() {
        if (this._timer) {
            clearInterval(this._timer);
            this._timer = null;
        }
        this._viewing = false;
        if (this._videoDevice) {
            this._videoDevice.stop();
        }
    },
    _capture: function() {
        var $control = this.element,
            id = $control.attr('id');

        this.turnOff();
        this._captureButton.text(this._originalButtonText);
        this.options.data = this._canvas.toDataURL("image/" + this.options.mimeType, this.options.quality);
        goradd.setControlValue(id, "data", this.options.data);
    },
    _loadmetadata: function() {
        var hCanvas = this._canvas.height;
        var wCanvas = this._canvas.width;
        var hScale = this._videoElement.videoHeight / hCanvas;
        var wScale = this._videoElement.videoWidth / wCanvas;
        var scale = this.options.zoom / 100 + 1;
        if (wScale < hScale) {
            this._w = scale * wCanvas;
            this._h = scale * this._videoElement.videoHeight / wScale ;
            this._x = 0;
            this._y = scale * (this._h - hCanvas) / 2 ;
        } else {
            this._h = scale * hCanvas;
            this._w = scale * this._videoElement.videoWidth / hScale;
            this._y = 0;
            this._x = scale * (this._w - wCanvas) / 2;

        }

    },
    _playing: function() {
        var self = this;

        if (!self._timer) {
            self._timer = setInterval(function () {
                try {
                    self._context.drawImage(self._videoElement, -self._x, -self._y, self._w, self._h);
                } catch (err) {
                    console.debug(err)
                }
            }, 100);
        }
    },
    _view: function() {
        var self = this;

        this._viewing = true;
        this._captureButton.text(this.options.selectButtonName);

        if (self._devices && self._devices.length > 1) {
            self._switchButton.show();
        } else {
            self._switchButton.hide();
        }

        var constraints = {video: true};

        if (this._devices) {
            constraints = {video: {devicedId: {exact: this._devices[this._deviceNum].devicedId}}}
        }
        navigator.mediaDevices.getUserMedia(constraints).then( function(mediaStream) {
            if (HTMLMediaElement) {
                self._videoElement.srcObject = mediaStream;
            } else {
                self._videoElement.src = URL.createObjectURL(mediaStream);
            }

            self._videoDevice = mediaStream.getVideoTracks()[0];
            self._videoElement.play(); // Safari bug


            var hCanvas = self._canvas.height;
            var wCanvas = self._canvas.width;

            self._context.restore();
            self._context.save();

            // flip image so selfie is easier to take
            var c = self._videoDevice.getCapabilities();

            if (!c.facingMode || c.facingMode === "user" || c.facingMode.length == 0 || c.facingMode[0] == "user") {
                self._context.translate(self._canvas.width, 0);
                self._context.scale(-1, 1);
            }

            if (self.options.shape == "circle") {
                self._context.beginPath();
                self._context.rect(0, 0, wCanvas, hCanvas);
                self._context.fillStyle = "white";
                self._context.fill();
                self._context.beginPath();
                var minDim = Math.min(hCanvas, wCanvas);
                self._context.arc(wCanvas / 2, hCanvas / 2, minDim / 2, 0, 2 * Math.PI, false);
                self._context.clip();
            }


        }).catch(function(err) {
            self.turnOff();
            self._showErr(err);
        });
    },
    _switchButtonClick: function() {
        if (this._deviceNum >= this._devices.length) {
            this._deviceNum = 0;
        } else {
            this._deviceNum ++;
        }
        this.turnOff();
        this._view();
    },
    _destroy: function () {
        this.turnOff();
        this._super();
    },
    _setOption: function( key, value ) {
        this._super( key, value );
        if ( key === "data" ) {
            this._loadNewData();
        }
    }
});
});
