package se.cag.jfocus.badgescanner.spring.hw;

import org.springframework.context.annotation.Profile;
import org.springframework.stereotype.Component;
import se.cag.jfocus.badgescanner.application.QrCodeRepo;
import se.cag.jfocus.badgescanner.domain.Player;

import java.util.Optional;
import java.util.UUID;

@Component
@Profile("local")
public class QrCodeReaderSim implements QrCodeRepo {

    private static int count = 0;
    @Override
    public Optional<Player> readUser() {
        count ++;
        if (count % 5 == 0) {
            String uuid = String.valueOf(UUID.randomUUID());
            try {
                Thread.sleep(3000);
            } catch (InterruptedException e) {
                throw new RuntimeException(e);
            }
            return Optional.of(Player.builder().email(uuid).name(uuid).company(uuid).build());
        }
        return Optional.empty();
    }
}
