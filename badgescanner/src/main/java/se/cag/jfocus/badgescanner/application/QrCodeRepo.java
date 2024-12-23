package se.cag.jfocus.badgescanner.application;

import se.cag.jfocus.badgescanner.annotations.JfocusRepo;
import se.cag.jfocus.badgescanner.domain.Player;

import java.util.Optional;

@JfocusRepo
public interface QrCodeRepo {

    Optional<Player> readUser();

}
