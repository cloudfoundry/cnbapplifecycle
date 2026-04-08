package src.main.java.org.cloudfoundry.test;

import org.cloudfoundry.test.core.InitializationUtils;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class Application {

    public static void main(String[] args) {
        new InitializationUtils().fail();
        SpringApplication.run(Application.class, args);
    }

}