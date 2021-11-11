import {AfterViewInit, ChangeDetectionStrategy, Component} from '@angular/core';

@Component({
    selector: 'app-networks',
    template: `
        <div class="mat-elevation-z8 content-elm">
            <h1>Add Network</h1>
            <form>
                <mat-grid-list [cols]="3" [rowHeight]="100">
                    <mat-grid-tile [colspan]="3">
                        <mat-form-field appearance="outline">
                            <mat-label>Name</mat-label>
                            <input matInput placeholder="Name">
                            <mat-hint>Hint</mat-hint>
                        </mat-form-field>
                    </mat-grid-tile>
                    <mat-grid-tile>
                        <mat-form-field appearance="outline">
                            <mat-label>CA Duration</mat-label>
                            <input matInput placeholder="CA Duration">
                            <mat-hint>Amount of time the certificate should be valid for</mat-hint>
                        </mat-form-field>
                    </mat-grid-tile>
                    <mat-grid-tile>
                        <mat-form-field appearance="outline">
                            <mat-label>Groups</mat-label>
                            <input matInput placeholder="Groups">
                            <mat-hint>This will limit which groups subordinate certs can use</mat-hint>
                        </mat-form-field>
                    </mat-grid-tile>
                    <mat-grid-tile>
                        <mat-form-field appearance="outline">
                            <mat-label>IP's</mat-label>
                            <input matInput placeholder="IP's">
                            <mat-hint>This will limit which ip addresses and networks subordinate certs can use
                            </mat-hint>
                        </mat-form-field>
                    </mat-grid-tile>
                    <mat-grid-tile>
                        <mat-form-field appearance="outline">
                            <mat-label>IP Pool</mat-label>
                            <input matInput placeholder="This will be used to assign IP's">
                            <mat-hint>Hint</mat-hint>
                        </mat-form-field>
                    </mat-grid-tile>
                    <mat-grid-tile>
                        <mat-form-field appearance="outline">
                            <mat-label>Subnet's</mat-label>
                            <input matInput placeholder="">
                            <mat-hint>This will limit which subnet addresses and networks subordinate certs can use
                            </mat-hint>
                        </mat-form-field>
                    </mat-grid-tile>
                </mat-grid-list>
            </form>
        </div>
        <div class="mat-elevation-z8 content-elm">
            <div style="justify-content: center; display: flex;">
                <button mat-raised-button color="primary" (click)="add()">Add</button>
            </div>
        </div>
    `,
    styles: [`
      table {
        width: 100%;
      }

      .content-elm {
        padding: 1em;
        margin-bottom: 1em;
      }

      mat-form-field {
        width: 100%;
        margin: 1em;
      }
    `],
    changeDetection: ChangeDetectionStrategy.OnPush,
})
export class NetworkAddComponent implements AfterViewInit {

    constructor() {

    }

    ngAfterViewInit() {

    }

    add() {

    }
}