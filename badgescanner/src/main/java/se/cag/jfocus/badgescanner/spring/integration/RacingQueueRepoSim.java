package se.cag.jfocus.badgescanner.spring.integration;

import lombok.extern.slf4j.Slf4j;
import org.springframework.context.annotation.Profile;
import org.springframework.stereotype.Component;
import se.cag.jfocus.badgescanner.application.RacingQueueRepo;
import se.cag.jfocus.badgescanner.domain.Player;

@Component
@Slf4j
@Profile("local")
public class RacingQueueRepoSim implements RacingQueueRepo {

    @Override
    public void enqueuePlayer(Player player) {
        log.info("Enqueue {}", player.name());
    }
}
