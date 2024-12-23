package se.cag.jfocus.badgescanner.spring.hw;

import com.google.zxing.BinaryBitmap;
import com.google.zxing.ChecksumException;
import com.google.zxing.FormatException;
import com.google.zxing.NotFoundException;
import com.google.zxing.RGBLuminanceSource;
import com.google.zxing.Result;
import com.google.zxing.common.HybridBinarizer;
import com.google.zxing.qrcode.QRCodeReader;
import com.pi4j.context.Context;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.context.annotation.Profile;
import org.springframework.stereotype.Component;
import se.cag.jfocus.badgescanner.application.QrCodeRepo;
import se.cag.jfocus.badgescanner.domain.Player;
import se.cag.jfocus.badgescanner.spring.BackendApplicationConfiguration;

import javax.imageio.ImageIO;
import java.awt.image.BufferedImage;
import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.util.Objects;
import java.util.Optional;


@Slf4j
@Component
@Profile("!local")
@RequiredArgsConstructor
public class QrCodeReader implements QrCodeRepo {

    private final BackendApplicationConfiguration configuration;
    private final Context pi4jContext;
    private final Lib4jCamera camera;
    private static final Lib4jCamera.PicConfig config = Lib4jCamera.newPictureConfigBuilder()
            .disablePreview(true)
            .timeout(1)
            .encoding(Lib4jCamera.PicEncoding.PNG)
            .quality(93)
            .width(820)
            .height(616)
            .build();


//    private static final CameraConfiguration camConfig = cameraConfiguration()
//            // 820x616 -- 1920x1080
//            .width(820)
//            .height(616)
//            .encoding(Encoding.JPEG)
//            .quality(85);

    public Optional<Player> readUser() {

        QrCodeDetector detector = new QrCodeDetector();
        camera.recordPicture(config, detector);
        if (Objects.nonNull(detector.result())) {
            return Optional.of(detector.result()).map(this::toRegisteredUser);
        }
        return Optional.empty();
    }

    private Player toRegisteredUser(Result result) {
        String data = result.getText();
        int firstIndex = data.indexOf(configuration.getSeparator());
        int lastIndex = data.lastIndexOf(configuration.getSeparator());
        return Player.builder()
                .name(firstIndex > 0 ? data.substring(0, firstIndex).trim() : "")
                .company(firstIndex > 0 && lastIndex > firstIndex ? data.substring(firstIndex + 1, lastIndex).trim(): "")
                .email(lastIndex > 0 ? data.substring(lastIndex + 1).trim() : "")
                .fullString(result.getText())
                .build();
    }

    private static class QrCodeDetector implements Lib4jCameraCallback {

        private Result result;

        @Override
        public void scannedImage(ByteArrayOutputStream imageData) {
            BinaryBitmap bitmap = null;
            try {
                BufferedImage image = ImageIO.read(new ByteArrayInputStream(imageData.toByteArray()));
                int[] pixels = image.getRGB(0, 0, image.getWidth(), image.getHeight(), null, 0, image.getWidth());
                RGBLuminanceSource source = new RGBLuminanceSource(image.getWidth(), image.getHeight(), pixels);
                bitmap = new BinaryBitmap(new HybridBinarizer(source));

            } catch (IOException e) {
                e.printStackTrace();
            }
            if (bitmap == null){
                log.error("Bitmap is null as the Image is not in a recognizable format, hence the Decode is unsuccessful!!");
            }
            QRCodeReader reader = new QRCodeReader();
            try {
                result = reader.decode(bitmap);
                log.info("Scan Decode is successful!!! The decoded QR Code text is:");
                log.info(result.getText());
                log.info("========End of processing the snap =======");
            } catch (NotFoundException e) {
                log.trace("Image is not in a recognizable format, hence the Decode is unsuccessful because of the following exception:", e);
                log.trace("This NotFoundException is thrown when a QR Code was not found in the image. It might have been partially detected but could not be confirmed.");
            } catch (ChecksumException e) {
                log.trace("Image is not in a recognizable format, hence the Decode is unsuccessful because of the following exception:", e);
                log.trace("This ChecksumException is thrown when a QR Code was successfully detected and decoded, but was not returned because its checksum feature failed.");
            } catch (FormatException e) {
                log.trace("Image is not in a recognizable format, hence the Decode is unsuccessful because of the following exception:", e);
                log.trace("This FormatException is thrown when a QR Code was successfully detected, but some aspect of the content did not conform to the barcode's format rules. This could have been due to a mis-detection.");
            }
            log.info("Image Processing Completed");
        }

        public Result result() {
            return result;
        }
    }

}
