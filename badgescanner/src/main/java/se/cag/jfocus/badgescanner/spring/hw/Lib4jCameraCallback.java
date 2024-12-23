package se.cag.jfocus.badgescanner.spring.hw;

import java.io.ByteArrayOutputStream;

public interface Lib4jCameraCallback {
    void scannedImage(ByteArrayOutputStream bytes);

}
