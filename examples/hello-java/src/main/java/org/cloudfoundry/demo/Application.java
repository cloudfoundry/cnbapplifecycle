package src.main.java.org.cloudfoundry.demo;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

@SpringBootApplication(proxyBeanMethods=false)
@RestController
public class Application {

	public static void main(String[] args) {
		SpringApplication.run(Application.class, args);
	}

	// CPU-intensive Fibonacci calculation
	private long fib(int n) {
		if (n <= 1) return n;
		return fib(n - 1) + fib(n - 2);
	}

	@GetMapping("/")
	public String index() {
		String instanceIndex = System.getenv("CF_INSTANCE_INDEX");
		
		// Do CPU-intensive work
		long fibResult = fib(40);
		System.out.println("Computed Fibonacci(40): " + fibResult);
		
		return "Hello World! (CF_INSTANCE_INDEX: " + instanceIndex + ")\nFibonacci(40): " + fibResult + "\n";
	}
}