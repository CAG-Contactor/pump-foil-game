package se.cag.jfocus.badgescanner.spring;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.ComponentScan;
import org.springframework.context.annotation.FilterType;
import org.springframework.scheduling.annotation.EnableScheduling;
import se.cag.jfocus.badgescanner.annotations.JfocusRepo;
import se.cag.jfocus.badgescanner.annotations.JfocusService;

@SpringBootApplication
@EnableScheduling
@ComponentScan(
        includeFilters = {
                @ComponentScan.Filter(type = FilterType.ANNOTATION, classes = { JfocusService.class, JfocusRepo.class}),
        },
        basePackages = {
                "se.cag.jfocus"
        }
)
public class BackendApplication {

    public static void main(String[] args) {
        SpringApplication.run(BackendApplication.class, args);
    }

}
