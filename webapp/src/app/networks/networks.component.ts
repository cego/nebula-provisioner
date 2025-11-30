import {Component, OnDestroy, OnInit} from '@angular/core';
import {Network} from "../models/network";
import {of as observableOf} from "rxjs";
import {catchError, map} from "rxjs/operators";
import {Apollo, gql, onlyCompleteData} from "apollo-angular";
import {ApolloResponse} from "../models/apollo";
import {SubSink} from "subsink";
import {Router} from "@angular/router";

@Component({
    selector: 'app-networks',
    template: `
        <div class="mat-elevation-z8">
            @if (isLoadingResults) {
                <div class="loading-shade">
                    @if (isLoadingResults) {
                        <mat-spinner></mat-spinner>
                    }
                </div>
            }

            <table mat-table [dataSource]="data">

                <ng-container matColumnDef="name">
                    <th mat-header-cell *matHeaderCellDef> Name</th>
                    <td mat-cell *matCellDef="let item"> {{ item.name }}</td>
                </ng-container>
                <ng-container matColumnDef="ips">
                    <th mat-header-cell *matHeaderCellDef> IP's</th>
                    <td mat-cell *matCellDef="let item"> {{ item.ips }}</td>
                </ng-container>
                <ng-container matColumnDef="ipPools">
                    <th mat-header-cell *matHeaderCellDef> IP Pools</th>
                    <td mat-cell *matCellDef="let item"> {{ item.ipPools }}</td>
                </ng-container>

                <tr mat-header-row *matHeaderRowDef="displayedColumns"></tr>
                <tr mat-row class="mat-row-selectable" *matRowDef="let row; columns: displayedColumns;"
                    (click)="openNetwork(row)"></tr>
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
    standalone: false
})
export class NetworksComponent implements OnInit, OnDestroy {
    private subs = new SubSink();
    displayedColumns: string[] = ['name', 'ips', 'ipPools'];
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
                        ips
                        ipPools
                    }
                }`,
        })
            .pipe(
                map(res => {
                    this.isLoadingResults = false;
                    if (res.data) {
                        return res.data.getNetworks;
                    } else {
                        throw new Error("Networks not found");
                    }
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
        this.router.navigate(['networks', network.name])
    }
}
