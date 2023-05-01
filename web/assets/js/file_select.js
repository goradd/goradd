(function(){
    goradd.FileSelect = class extends goradd.Widget {
        constructor(element, options) {
            let optionDefaults = {

            };
            options = goradd.extendOptions(optionDefaults, options);
            super(element, options);
            this._files = [];
            this._readyToUpload = false;
            this.on('change', this._handleChange, {bubbles: false});

        }
        _handleChange (event) {
            this._files = event.target.files;
        }
        val(v) {
            if (this._readyToUpload) {
                this._readyToUpload = false;
                return this._files;
            } else {
                var ret = []

                // returns a stringified array about all the files
                for (var i = 0; i < this._files.length; i++) {
                    var f = this._files[i];
                    var o = {
                        lastModified: f.lastModified,
                        name: f.name,
                        size: f.size,
                        type: f.type
                    }
                    ret.push(o);
                }
                return JSON.stringify(ret);

            }
        }
        upload() {
            this._readyToUpload = true;
            this.trigger("formObjChanged");
            this.trigger("upload");
        }
    };

    goradd.registerWidget("goradd.FileSelect", goradd.FileSelect);

})();