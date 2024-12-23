package se.cag.jfocus.badgescanner.spring.integration;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.context.annotation.Profile;
import org.springframework.http.RequestEntity;
import org.springframework.stereotype.Component;
import org.springframework.web.client.HttpStatusCodeException;
import org.springframework.web.client.RestTemplate;
import se.cag.jfocus.badgescanner.application.RacingQueueRepo;
import se.cag.jfocus.badgescanner.domain.Player;
import se.cag.jfocus.badgescanner.spring.BackendApplicationConfiguration;
import se.cag.jfocus.badgescanner.spring.integration.dto.NewUser;

import java.net.URI;

@Component
@RequiredArgsConstructor
@Slf4j
@Profile("!local")
public class RacingQueueRepoApi implements RacingQueueRepo {

    private final RestTemplate template;
    private final BackendApplicationConfiguration configuration;

    @Override
    public void enqueuePlayer(Player player) {
        log.info("Enqueue {}", player.name());
        RequestEntity<NewUser> body = RequestEntity.post(URI.create(configuration.getEnqueueUrl()))
                .body(toNewUser(player));
        try {
            template.exchange(body, Void.class);
        } catch (HttpStatusCodeException e) {
            log.error("Failed to enqueue player: " + e.getMessage());
        }
    }

    private NewUser toNewUser(Player player) {
        return NewUser.builder()
                .userId("fakeId")
                .displayName(player.name())
                .organisation(player.company())
                .password("112233")
                .build();
    }
}
