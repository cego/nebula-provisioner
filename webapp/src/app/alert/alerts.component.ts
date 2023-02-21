import {Component, OnDestroy, OnInit, ViewEncapsulation} from '@angular/core';
import {SubSink} from "subsink";
import {Alert, AlertService} from "./alert.service";
import {MatLegacySnackBar as MatSnackBar} from "@angular/material/legacy-snack-bar";


@Component({
    template: ``,
    selector: 'app-alerts',
    styles: [`

      .snack-bar-danger {
        background-color: darkred;
        max-width: 900px
      }
    `],
    encapsulation: ViewEncapsulation.None,
})
export class AlertsComponent implements OnInit, OnDestroy {
    private subs = new SubSink();

    constructor(private alertService: AlertService, private snackBar: MatSnackBar) {
    }

    ngOnInit(): void {
        this.subs.add(this.alertService.$alerts
            .subscribe(alert => {
                this.addAlert(alert);
            }))
    }

    ngOnDestroy(): void {
        this.subs.unsubscribe();
    }

    addAlert(alert: Alert) {
        this.snackBar.open(alert.message, 'Close', {
            duration: (alert.duration) ? alert.duration : 5 * 1000,
            horizontalPosition: 'center',
            panelClass: [`snack-bar`, `snack-bar-${alert.type}`],
        });
    }
}
