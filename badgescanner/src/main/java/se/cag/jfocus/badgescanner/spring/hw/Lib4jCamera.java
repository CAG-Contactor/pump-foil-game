package se.cag.jfocus.badgescanner.spring.hw;

import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Component;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.InputStream;

@Component
@Slf4j
public class Lib4jCamera {

    /**
     * FHNW implementation of a camera, works with the raspberry-pi v2 camera module and
     * the Pi4J-Basic-OS image on the raspberry-pi.
     * <p>
     * Maybe works on other camera-modules too, but is not yet tested.
     * <p>
     * It uses the libcamera-still and libcamera-vid bash commands. those are pre-installed
     * on all raspbian-versions after Buster.
     */
    public static PicConfig.Builder newPictureConfigBuilder() {
        return new PicConfig.Builder();
    }

    /**
     * Constructor for using the picture and video functionality
     * calling the init function to test if a camera is active
     */
    public Lib4jCamera() {
        init();
    }

    /**
     * Takes a picture and saves it to the default Pictures folder
     * <p>
     * If a file already exists, the code will break. better use useDate while taking pictures
     */

    /**
     * Takes a picture using the bash commands
     *
     * @param config Use the ConfigBuilder of the CameraConfig to create the desired parameters
     */
    public void recordPicture(PicConfig config, Lib4jCameraCallback callback) {
        log.debug("Taking Picture");

        ProcessBuilder processBuilder = new ProcessBuilder();
        processBuilder.command("bash", "-c", config.asCommand());

        try {
            callBash(processBuilder, callback);
        } catch (Exception e) {
            log.error("Camera: Error while taking picture: ", e);
        }
    }

    /**
     * Uses a ProcessBuilder to call the bash of the RaspberryPI.
     * This will call the command and write the output to the console
     *
     * @param processBuilder which process needs to be built
     */
    private void callBash(ProcessBuilder processBuilder, Lib4jCameraCallback callback) throws IOException, InterruptedException {
        Process process = processBuilder.start();
        ByteArrayOutputStream image = new ByteArrayOutputStream();

        int len = -1;
        try (InputStream stream = process.getInputStream()) {
            image.writeBytes(stream.readAllBytes());

            //exitCode 0 = No Errors
            int exitCode = process.waitFor();

            if (exitCode != 0) {
                log.error("Camera exited with error code : %s", exitCode);
            } else {
                callback.scannedImage(image);
                log.info("Camera finished successfully");
            }
        }
    }

    private void callBash(ProcessBuilder processBuilder) throws IOException, InterruptedException {
        Process process = processBuilder.start();
        //exitCode 0 = No Errors0
        int exitCode = process.waitFor();
        if (exitCode != 0) {
            log.error("Camera exited with error code : %s", exitCode);
        } else {
            log.info("Camera finished successfully");
        }
    }

    /**
     * testing, if camera is installed on raspberrypi, and if the bash commands
     * will work
     */
    private void init() {
        log.debug("initialisation of camera");

        ProcessBuilder processBuilder = new ProcessBuilder();
        processBuilder.command("bash", "-c", "libcamera-still");

        try {
            callBash(processBuilder);
        } catch (Exception e) {
            log.error("Camera: Error at initialisation: ", e);
        }
    }

    /**
     * Output Format of pictures
     * These modes determine the output of the picture-file
     * <p>
     * The following encodings can be set
     * {@link #PNG}
     * {@link #JPG}
     * {@link #RGB}
     * {@link #BMP}
     * {@link #YUV420}
     */
    public enum PicEncoding {
        PNG("png"),
        JPG("jpg"),
        RGB("rgb"),
        BMP("bmp"),
        YUV420("yuv420");

        private final String encoding;

        PicEncoding(String encoding) {
            this.encoding = encoding;
        }

        public String getEncoding() {
            return encoding;
        }
    }


    /**
     * Builder Pattern to create a config for a single Picture
     */
    public static class PicConfig {
         /**
         * a delay, before taking a picture
         */
        public final int delay;

        public final int timeout;
        /**
         * output width of the picture
         */
        public final int width;
        /**
         * output height of the picture
         */
        public final int height;
        /**
         * the quality of the picture, ranging from 0 to 100
         * where 100 is the best quality of the picture, with no blurring
         */
        public final int quality;
        /**
         * The format of the output
         */
        public final PicEncoding encoding;
        /**
         * when true, there is no preview on the raspberry-pi
         */
        public final boolean disablePreview;
        /**
         * when true, the preview is in fullscreen
         */
        public final boolean allowFullscreenPreview;

        /**
         * constructor for the config
         *
         * @param builder builder with the defined options
         */
        private PicConfig(Builder builder) {
            this.delay = builder.delay;
            this.width = builder.width;
            this.height = builder.height;
            this.quality = builder.quality;
            this.encoding = builder.encoding;
            this.timeout = builder.timeout;
            this.disablePreview = builder.disablePreview;
            this.allowFullscreenPreview = builder.allowFullscreenPreview;
        }

        /**
         * Creates a callable bash command with the defined options.
         *
         * @return a string that can be called from the bash
         */
        public String asCommand() {
            StringBuilder command = new StringBuilder("libcamera-still");
            command.append(" -o - ");
            if (delay != 0) {
                command.append(" -t ").append(delay);
            }
            if (width != 0) {
                command.append(" --width ").append(width);
            }
            if (height != 0) {
                command.append(" --height ").append(height);
            }
            if (quality != 0) {
                command.append(" -q ").append(quality);
            }
            if (encoding != null) {
                command.append(" --encoding ").append(encoding.getEncoding());
            }
            if (disablePreview) {
                command.append(" -n");
            }
            if (allowFullscreenPreview && !disablePreview) {
                command.append(" -f");
            }
            if (timeout != 0) {
                command.append(" -t ").append(timeout);
            }

            return command.toString();
        }

        /**
         * Builder Pattern, to create a config for a single picture
         * <p>
         * A Config is buildable like this:
         * var config = Camera.PicConfig.Builder.newInstance()
         * .outputPath("/home/pi/Pictures/")
         * .delay(3000)
         * .disablePreview(true)
         * .encoding(Camera.PicEncoding.PNG)
         * .useDate(true)
         * .quality(93)
         * .width(1280)
         * .height(800)
         * .build();
         * <p>
         * Every property can be added or not.
         */
        public static class Builder {
            private int delay;
            private int timeout;
            private int width;
            private int height;
            private int quality;
            private PicEncoding encoding;
            private boolean disablePreview;
            private boolean allowFullscreenPreview;


            public Builder delay(int delay) {
                this.delay = delay;
                return this;
            }

            public Builder timeout(int timeout) {
                this.timeout = timeout;
                return this;
            }

            public Builder width(int width) {
                this.width = width;
                return this;
            }

            public Builder height(int height) {
                this.height = height;
                return this;
            }

            public Builder quality(int quality) {
                if (quality < 0 || quality > 100) {
                    throw new IllegalArgumentException("quality must be between 0 and 100");
                }
                this.quality = quality;
                return this;
            }

            public Builder encoding(PicEncoding encoding) {
                this.encoding = encoding;
                return this;
            }

            public Builder disablePreview(boolean disablePreview) {
                this.disablePreview = disablePreview;
                return this;
            }

            public Builder allowFullscreenPreview(boolean allowFullscreenPreview) {
                this.allowFullscreenPreview = allowFullscreenPreview;
                return this;
            }

            public PicConfig build() {
                return new PicConfig(this);
            }
        }
    }

}
