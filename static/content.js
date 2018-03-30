            $.fn.loadImages = function (callback) {
                var element = this;
                $.getJSON("http://localhost:8081/files", function (data, textStatus) {
                    $.each(data, function (i, item) {
                        var title = "pic";
                        var link = "http://localhost:8081/files/" + item.user + "/" + item.name;
                        var thumbnail = "http://localhost:8081/files/" + item.user + "/thumbnail/" + item.name;
                        $("<a/>").attr({
                            "href": link
                        }).append(
                                $("<img/>").attr({
                                    "src": thumbnail,
                                    "alt": title,
                                    "title": title
                                }).css({
                                    "margin": "2px",
                                    "border": "none",
                                    "vertical-align": "bottom"
                                })
                        ).appendTo(element);
                    });
                    callback();
                });
            };
            $(document).ready(function () {
                $("#images").loadImages(function () {
                    $("#images a").visage();
                });
                setTimeout(function () {
                    $("#images-alt").loadImages(function () {
                        $("#images-alt a").visage({
                            "files": {
                                "blank": "/static/img/blank.gif",
                                "error": "/static/img/error.png"
                            },
                            "attr": {
                                "close": {"id": "visage-alt-close"},
                                "title": {"id": "visage-alt-title"},
                                "count": {"id": "visage-alt-count"},
                                "container": {"id": "visage-alt-container"},
                                "image": {"id": "visage-alt-image", "src": "/static/img/blank.gif"},
                                "visage": {"id": "visage-alt"},
                                "overlay": {"id": "visage-alt-overlay"},
                                "prev": {"id": "visage-alt-nav-prev"},
                                "prev_disabled": {"id": "visage-alt-nav-prev"},
                                "next": {"id": "visage-alt-nav-next"},
                                "next_disabled": {"id": "visage-alt-nav-next"}
                            },
                            "enabledClass": "visage-alt-enabled",
                            "disabledNavClass": "visage-alt-nav-disabled",
                            "transitionResizeSpeed": 300, // Non-zero to show resize animation
                            "addDOM": function (visageDOM, options) {
                                $.fn.visage.addDOM(visageDOM, options);
                                // Moves elements to overlay so they are all in the same stacking context
                                $(visageDOM.prev).add(visageDOM.next).add(visageDOM.count).add(visageDOM.title).appendTo(visageDOM.overlay);
                            },
                            // We move setting image source to preTransitionResize so that transition resize is also resizing the image
                            // (So that image is shown as soon as possible and we are not waiting for the resize to finish to display it)
                            "preTransitionResize": function (image, values, group, index, visageDOM, isStopping, finish, options) {
                                $.fn.visage.preTransitionResize(image, values, group, index, visageDOM, isStopping, function () {}, options);
                                visageDOM.image.attr("src", values.src);
                                finish();
                            },
                            "postTransitionResize": function (image, values, group, index, visageDOM, isStopping, finish, options) {
                                finish();
                            }
                        });
                    });
                }, 2000);
            });
            Dropzone.options.dropzoneId = {
                paramName: "file",
                maxFileSie: 10,
                acceptedFiles: "image/*",

                success: function (file, response) {
                    if (response.status == 406) {
                        alert("File type not supported");
                        return;
                    }
                    fileinfo = JSON.parse(response);
                    var fileuploded = file.previewElement.querySelector("[data-dz-name]");
                    fileuploded.innerHTML = fileinfo.name;
                }

            }
