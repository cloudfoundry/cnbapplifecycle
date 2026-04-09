package src.main.java.org.cloudfoundry.demo;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.ArrayList;
import java.util.List;

@SpringBootApplication(proxyBeanMethods=false)
@RestController
public class Application {

	public static void main(String[] args) {
		SpringApplication.run(Application.class, args);
	}

	@GetMapping("/")
	public String fibonacci() {
		List<Long> fibNumbers = new ArrayList<>();
		long a = 0, b = 1;
		
		for (int i = 0; i < 40; i++) {
			fibNumbers.add(a);
			long next = a + b;
			a = b;
			b = next;
		}
		
		return fibNumbers.toString();
	}

}