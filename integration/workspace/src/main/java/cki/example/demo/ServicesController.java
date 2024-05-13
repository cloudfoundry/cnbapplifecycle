package cki.example.demo;

import com.sap.cloud.environment.servicebinding.api.DefaultServiceBindingAccessor;
import com.sap.cloud.environment.servicebinding.api.ServiceBinding;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;

@RestController
public class ServicesController {

    @GetMapping("/services")
    public String getServices(){
        StringBuilder sb = new StringBuilder();
        for (ServiceBinding svc : DefaultServiceBindingAccessor.getInstance().getServiceBindings()) {
            sb.append("Service instance name: ");
            sb.append(svc.getName().orElse("missing"));
            sb.append(", service plan: ");
            sb.append(svc.getServicePlan().orElse("missing"));
            sb.append("\n");
        }
        return sb.toString();
    }
}
