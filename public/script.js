(function () {
  const width = 480;
  let height = 0;

  const video = document.querySelector(".camera-feed");
  const canvas = document.querySelector(".camera-canvas");
  const list = document.querySelector(".objectList");

  // get the users webcam
  navigator.mediaDevices
    .getUserMedia({ video: true, audio: false })
    .then((stream) => {
      video.srcObject = stream;
      video.play();
    })
    .catch((err) => {
      console.error(err);
    });

  video.addEventListener("canplay", () => {
    height = video.videoHeight / (video.videoWidth / width);
    video.setAttribute("width", width);
    video.setAttribute("height", height);

    // now that the camera is playing back in a video element for us to see
    // let's grab a snapshot every 5 seconds
    setInterval(() => {
      grabPhoto();
    }, 5000);
  });

  function grabPhoto() {
    // draw the frame from the video element
    // into the canvas
    let context = canvas.getContext("2d");
    canvas.width = width;
    canvas.height = height;
    context.drawImage(video, 0, 0, width, height);
    // grab the image into base64 and strip the header from it
    let data = canvas.toDataURL("image/png").substring(22);

    // now lets send the image to our Go endpoint and return the data!
    fetch("http://localhost:3000/check", {
      method: "post",
      body: JSON.stringify({
        Image: data,
      }),
    }).then((res) => {
      res.json().then((content) => {
        // clear the current list
        list.innerHTML = "";
        const labels = content.Labels;

        for (let i = 0; i < labels.length; i++) {
          list.innerHTML += `<li>
              ${labels[i].Name} - ${labels[i].Confidence}
          </li>`;
        }
      });
    });
  }
})();
