package se.cag.jfocus.badgescanner.spring.control;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Component;
import se.cag.jfocus.badgescanner.application.RacingQueueService;

import java.util.concurrent.TimeUnit;

@Component
@RequiredArgsConstructor
@Slf4j
public class RacingQueueScheduler {

    private final RacingQueueService racingQueueService;

    @Scheduled(fixedDelay = 3, timeUnit = TimeUnit.SECONDS)
    public void enqueue() {
        log.info("Checking for new User");
        racingQueueService.enqueueNewUserByQr();
    }
}
