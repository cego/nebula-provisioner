import {NgModule} from '@angular/core';
import {BrowserModule} from '@angular/platform-browser';

import {AppRoutingModule} from './app-routing.module';
import {AppComponent} from './app.component';
import {GraphQLModule} from './graphql.module';
import {HttpClientModule} from '@angular/common/http';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {MatLegacySliderModule as MatSliderModule} from "@angular/material/legacy-slider";
import {MatToolbarModule} from "@angular/material/toolbar";
import {MatIconModule} from "@angular/material/icon";
import {MatLegacyButtonModule as MatButtonModule} from "@angular/material/legacy-button";
import {MatSidenavModule} from "@angular/material/sidenav";
import {MatLegacyMenuModule as MatMenuModule} from "@angular/material/legacy-menu";
import {MatLegacyListModule as MatListModule} from "@angular/material/legacy-list";
import {NetworksComponent} from './networks/networks.component';
import {UserApproveDialog, UserDeleteDialog, UsersComponent} from './users/users.component';
import {MatLegacyPaginatorModule as MatPaginatorModule} from "@angular/material/legacy-paginator";
import {MatLegacyTableModule as MatTableModule} from "@angular/material/legacy-table";
import {MatLegacyProgressSpinnerModule as MatProgressSpinnerModule} from "@angular/material/legacy-progress-spinner";
import {MatSortModule} from "@angular/material/sort";
import {NetworkAddComponent} from "./networks/network-add.component";
import {MatLegacyFormFieldModule as MatFormFieldModule} from "@angular/material/legacy-form-field";
import {MatLegacyInputModule as MatInputModule} from "@angular/material/legacy-input";
import {MatGridListModule} from "@angular/material/grid-list";
import {MatLegacyDialogModule as MatDialogModule} from "@angular/material/legacy-dialog";
import {NetworkComponent} from "./networks/network.component";
import {NetworkAgentRevokeDialog, NetworkAgentsComponent} from "./networks/network-agents.component";
import {NetworkEnrollmentRequestsComponent} from "./networks/network-enrollment-requests.component";
import {AlertsComponent} from "./alert/alerts.component";
import {MAT_LEGACY_SNACK_BAR_DEFAULT_OPTIONS as MAT_SNACK_BAR_DEFAULT_OPTIONS, MatLegacySnackBarModule as MatSnackBarModule} from "@angular/material/legacy-snack-bar";
import {EnrollmentRequestApproveDialog} from "./networks/enrollment-request-approve-dialog.component";
import {CdkTableModule} from "@angular/cdk/table";

@NgModule({
    declarations: [
        AppComponent,
        AlertsComponent,
        EnrollmentRequestApproveDialog,
        NetworksComponent,
        NetworkComponent,
        NetworkAddComponent,
        NetworkAgentsComponent,
        NetworkAgentRevokeDialog,
        NetworkEnrollmentRequestsComponent,
        UsersComponent,
        UserApproveDialog,
        UserDeleteDialog
    ],
    imports: [
        BrowserModule,
        AppRoutingModule,
        GraphQLModule,
        HttpClientModule,
        BrowserAnimationsModule,
        MatSliderModule,
        MatToolbarModule,
        MatIconModule,
        MatButtonModule,
        MatSidenavModule,
        MatMenuModule,
        MatListModule,
        MatPaginatorModule,
        MatTableModule,
        MatProgressSpinnerModule,
        MatSortModule,
        MatFormFieldModule,
        MatInputModule,
        MatGridListModule,
        MatDialogModule,
        MatSnackBarModule,
        CdkTableModule
    ],
    providers: [
        {provide: MAT_SNACK_BAR_DEFAULT_OPTIONS, useValue: {duration: 2500}}
    ],
    bootstrap: [AppComponent]
})
export class AppModule {
}
