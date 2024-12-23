package se.cag.jfocus.badgescanner.spring;

import com.pi4j.Pi4J;
import com.pi4j.context.Context;
import lombok.Data;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.boot.web.client.RestTemplateBuilder;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.client.RestTemplate;

@Data
@Configuration
@ConfigurationProperties(prefix = "cag.qr")
public class BackendApplicationConfiguration {
    private String separator;
    private String enqueueUrl;

    @Bean
    public RestTemplate template(RestTemplateBuilder builder) {
        return builder.build();
    }

    @Bean
    public Context createContext() {
        return Pi4J.newAutoContext();
    }

}
