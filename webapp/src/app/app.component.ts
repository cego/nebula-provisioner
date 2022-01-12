import {Component} from '@angular/core';

@Component({
    selector: 'app-root',
    template: `
        <p>
            <mat-toolbar color="primary">
                <span style="margin-right: 1em;">Nebula Provisioner</span>
                <a mat-button [routerLink]="['networks']">Networks</a>
                <a mat-button [routerLink]="['users']">Users</a>
                <span class="toolbar-spacer"></span>
            </mat-toolbar>
        </p>

        <div class="content">
            <router-outlet></router-outlet>
            <app-alerts></app-alerts>
        </div>
    `,
    styles: [`.toolbar-spacer {
      flex: 1 1 auto;
    }

    .content {
      margin: 1em auto;
      max-width: 100em;
      padding: 1em;
    }
    `],
})
export class AppComponent {
}