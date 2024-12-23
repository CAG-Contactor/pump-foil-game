package se.cag.jfocus.badgescanner.application;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import se.cag.jfocus.badgescanner.annotations.JfocusService;
import se.cag.jfocus.badgescanner.domain.Player;

import java.util.Optional;

@Slf4j
@JfocusService
@RequiredArgsConstructor
public class ScanPlayerService {

    private final QrCodeRepo repo;

    public Optional<Player> scanPlayer() {
       Optional<Player> player = repo.readUser();
       if (player.isPresent() && player.get().isValid()) {
           return player;
       }

        return Optional.empty();
    }

}
