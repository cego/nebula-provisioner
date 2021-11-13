import {Component, OnDestroy, OnInit} from '@angular/core';
import {Network} from "../models/network";
import {of as observableOf} from "rxjs";
import {catchError, map} from "rxjs/operators";
import {Apollo, gql} from "apollo-angular";
import {ApolloResponse} from "../models/apollo";
import {SubSink} from "subsink";
import {Router} from "@angular/router";

@Component({
    selector: 'app-networks',
    template: `
        <!--<div class="mat-elevation-z8 action-buttons">
            <div>
                <button mat-raised-button color="primary" [routerLink]="['add']">Add</button>
            </div>
        </div>-->

        <div class="mat-elevation-z8">
            <div class="loading-shade" *ngIf="isLoadingResults">
                <mat-spinner *ngIf="isLoadingResults"></mat-spinner>
            </div>

            <table mat-table [dataSource]="data">

                <!-- Name Column -->
                <ng-container matColumnDef="name">
                    <th mat-header-cell *matHeaderCellDef> Name</th>
                    <td mat-cell *matCellDef="let item"> {{item.name}} </td>
                </ng-container>
                <!-- IP pools Column -->
                <ng-container matColumnDef="ipPools">
                    <th mat-header-cell *matHeaderCellDef> IP Pools</th>
                    <td mat-cell *matCellDef="let item"> {{item.ipPools}} </td>
                </ng-container>

                <tr mat-header-row *matHeaderRowDef="displayedColumns"></tr>
                <tr mat-row *matRowDef="let row; columns: displayedColumns;" (click)="openNetwork(row)"></tr>
            </table>
        </div>
    `,
    styles: [`
      table {
        width: 100%;
      }

      .action-buttons {
        padding: 1em;
        margin-bottom: 1em;
      }
    `],
})
export class NetworksComponent implements OnInit, OnDestroy {
    private subs = new SubSink();
    displayedColumns: string[] = ['name', 'ipPools'];
    data: Network[] = [];

    isLoadingResults = false;

    constructor(private apollo: Apollo, private router: Router) {
    }

    ngOnInit(): void {
        this.isLoadingResults = true;

        this.subs.sink = this.apollo.query<ApolloResponse>({
            query: gql`
                {
                    getNetworks {
                        name
                        ipPools
                    }
                }`,
        })
        .pipe(map(res => {
                this.isLoadingResults = false;
                return res.data.getNetworks;
            }),
            catchError(() => observableOf([])))
        .subscribe(data => {
            this.data = data;
        });

    }

    ngOnDestroy(): void {
        this.subs.unsubscribe();
    }

    openNetwork(network: Network) {
        console.log(network);
        this.router.navigate(['networks', network.name])
    }
}