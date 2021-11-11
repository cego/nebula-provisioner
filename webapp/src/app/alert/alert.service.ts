import {Injectable} from "@angular/core";
import {Subject} from "rxjs";

export class Alert {
    type!: string;
    message!: string;
    duration?: number
}


@Injectable({
    providedIn: 'root'
})
export class AlertService {

    private sub = new Subject<Alert>()

    $alerts = this.sub.asObservable();

    add(alert: Alert) {
        this.sub.next(alert);
    }

    addAlert(type: string, message: string, duration?: number) {
        const alert = new Alert();
        alert.type = type;
        alert.message = message;
        if (duration) {
            alert.duration = duration;
        }

        this.add(alert);
    }
}
