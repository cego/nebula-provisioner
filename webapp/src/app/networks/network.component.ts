import {Component, OnDestroy, OnInit} from '@angular/core';
import {EnrollmentRequest, GET_NETWORK_BY_NAME, Network} from "../models/network";
import {Apollo} from "apollo-angular";
import {SubSink} from "subsink";
import {ActivatedRoute} from "@angular/router";
import {filter, switchMap} from "rxjs/operators";
import {ApolloResponse} from "../models/apollo";
import {Agent} from "../models/agent";
import {AlertService} from "../alert/alert.service";

@Component({
    selector: 'app-network',
    template: `
        <div class="mat-elevation-z8 content-elm">
            <h1>Network: {{network?.name}}</h1>
            <table>
                <tbody>
                <tr>
                    <td>IP Pools:</td>
                    <td>{{network?.ipPools}}</td>
                </tr>
                <tr>
                    <td>Duration:</td>
                    <td>{{network?.duration}}</td>
                </tr>
                <tr>
                    <td>Groups:</td>
                    <td>{{network?.groups}}</td>
                </tr>
                <tr>
                    <td>IP's:</td>
                    <td>{{network?.ips}}</td>
                </tr>
                <tr>
                    <td>Subnets:</td>
                    <td>{{network?.subnets}}</td>
                </tr>
                <tr>
                    <td>Enrollment token:</td>
                    <td>{{network?.enrollmentToken}}</td>
                </tr>
                </tbody>
            </table>
        </div>
        <div class="mat-elevation-z8 content-elm" *ngIf="enrollmentRequests.length > 0">
            <h1>Enrollment Requests</h1>

            <app-network-enrollment-requests
                    [enrollmentRequests]="enrollmentRequests"></app-network-enrollment-requests>
        </div>

        <div class="mat-elevation-z8 content-elm" *ngIf="agents.length > 0">
            <h1>Agents</h1>

            <app-network-agents [agents]="agents"></app-network-agents>
        </div>
    `,
    styles: [`
      table {
        width: 100%;
      }

      mat-form-field {
        width: 100%;
        margin: 1em;
      }

      .mat-column-actions {
        padding-top: 0.5em;
        padding-bottom: 0.5em;
        width: 10em;
        text-align: right;
      }
    `],
})
export class NetworkComponent implements OnInit, OnDestroy {

    private subs = new SubSink();

    network!: Network;
    enrollmentRequests: EnrollmentRequest[] = [];
    agents: Agent[] = [];

    constructor(private apollo: Apollo, private route: ActivatedRoute, private alert: AlertService) {

    }

    ngOnInit(): void {

        this.subs.sink = this.route.params
            .pipe(
                filter(p => p && p?.name !== ''),
                switchMap(params =>
                    this.apollo.watchQuery<ApolloResponse>({
                        variables: {
                            name: params.name
                        },
                        query: GET_NETWORK_BY_NAME,
                    }).valueChanges
                )
            )
            .subscribe(data => {
                this.network = data?.data?.getNetwork;
                this.enrollmentRequests = (this.network.enrollmentRequests) ? this.network.enrollmentRequests : [];
                this.agents = (this.network.agents) ? this.network.agents : [];
            }, error => {
                this.alert.addAlert('danger', error.message);
            });
    }

    ngOnDestroy(): void {
        this.subs.unsubscribe();
    }


}