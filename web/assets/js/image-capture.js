/**
 * ImageCapture is the javascript support for the ImageCapture widget.
 * Since the image capture API is only supported by ES6 capable browsers, we define the widget using ES6 style classes.
 */
(function(){
    goradd.ImageCapture = class extends goradd.Widget {
        constructor(element, options) {
            let optionDefaults = {
                selectButtonName: "Capture",          // Name that the button gets when in capture mode
                shape: "rect",             // Mask shape
                zoom: 0,                  // % of zoom. 100 means 2x zoom. 200 means 3x zoom.
                data: "data:image/jpeg;base64,",   // Mime typed and base64'd data.
                mimeType: "jpeg",            // desired mime type
                quality: 0.92               // Quality setting for jpeg and webP
            };
            options = goradd.extendOptions(optionDefaults, options);
            super(element, options);
            this._viewing = false;
            if  (!('mediaDevices' in navigator)) {
                // either browser or device does not support media capture
                this._showErr();
                return;
            }
            this._videoElement = goradd.tagBuilder('video').attr('playsinline', "").element();
            this._canvas = this.find("canvas");
            this._context = this._canvas.element.getContext('2d');
            this._captureButton = g$(this.id + "-capture");
            this._switchButton = g$(this.id + "-switch");
            this._originalButtonText = this._captureButton.text();
            this._captureButton.on("click", [this, this._imageButtonClick]);
            this._switchButton.on("click", [this, this._switchButtonClick]);
            g$(this._videoElement).on("playing", [this, this._playing]);
            g$(this._videoElement).on("loadedmetadata", [this, this._loadmetadata]);
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
        }

        _showErr(err) {
            this._captureButton.prop("disabled", true);
            if (err) {
                let p = g$(this.id + "-err");
                p.show();
                p.text(err);
            }
        }
        async _loadNewData() {
            if (this.options && this.options.data) {
                let self = this;
                let p = new Promise((resolve, reject) => {
                    let image = new Image();
                    image.onload = () => resolve(image);
                    image.src = self.options.data;
                });
                let img = await p;
                this._context.drawImage(img, 0, 0);
            }
        }
        _imageButtonClick() {
            if (this._viewing) {
                this._capture();
            } else {
                this._view();
            }
        }
        turnOff() {
            if (this._timer) {
                clearInterval(this._timer);
                this._timer = null;
            }
            this._viewing = false;
            if (this._videoDevice) {
                this._videoDevice.stop();
            }
        }
        _capture() {
            this._captureButton.text(this._originalButtonText);
            this.options.data = this._canvas.element.toDataURL("image/" + this.options.mimeType, this.options.quality);
            goradd.setControlValue(this.id, "data", this.options.data);
            this.turnOff();
            this.trigger("capture");
        }
        _loadmetadata() {
            let hCanvas = this._canvas.element.height;
            let wCanvas = this._canvas.element.width;
            let hScale = this._videoElement.videoHeight / hCanvas;
            let wScale = this._videoElement.videoWidth / wCanvas;
            let scale = this.options.zoom / 100 + 1;
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

        }
        _playing() {
            let self = this;

            if (!self._timer) {
                self._timer = setInterval(function () {
                    try {
                        self._context.drawImage(self._videoElement, -self._x, -self._y, self._w, self._h);
                    } catch (err) {
                        console.debug(err)
                    }
                }, 100);
            }
        }
        _view() {
            let self = this;

            this._viewing = true;
            this._captureButton.text(this.options.selectButtonName);

            if (self._devices && self._devices.length > 1) {
                self._switchButton.show();
            } else {
                self._switchButton.hide();
            }

            let constraints = {video: true};

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
                self._videoElement.play();

                var hCanvas = self._canvas.element.height;
                var wCanvas = self._canvas.element.width;

                if (self.options.shape === "circle") {
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
        }
        _switchButtonClick() {
            if (this._deviceNum >= this._devices.length) {
                this._deviceNum = 0;
            } else {
                this._deviceNum ++;
            }
            this.turnOff();
            this._view();
        }
        _destroy() {
            this.turnOff();
            this._super();
        }
        _setOption( key, value ) {
            super._setOption( key, value );
            if ( key === "data" ) {
                this._loadNewData();
            }
        }

    };

    goradd.registerWidget("goradd.ImageCapture", goradd.ImageCapture);

})();
