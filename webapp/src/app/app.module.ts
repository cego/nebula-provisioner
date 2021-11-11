import {NgModule} from '@angular/core';
import {BrowserModule} from '@angular/platform-browser';

import {AppRoutingModule} from './app-routing.module';
import {AppComponent} from './app.component';
import {GraphQLModule} from './graphql.module';
import {HttpClientModule} from '@angular/common/http';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {MatSliderModule} from "@angular/material/slider";
import {MatToolbarModule} from "@angular/material/toolbar";
import {MatIconModule} from "@angular/material/icon";
import {MatButtonModule} from "@angular/material/button";
import {MatSidenavModule} from "@angular/material/sidenav";
import {MatMenuModule} from "@angular/material/menu";
import {MatListModule} from "@angular/material/list";
import {NetworksComponent} from './networks/networks.component';
import {UserApproveDialog, UserDeleteDialog, UsersComponent} from './users/users.component';
import {MatPaginatorModule} from "@angular/material/paginator";
import {MatTableModule} from "@angular/material/table";
import {MatProgressSpinnerModule} from "@angular/material/progress-spinner";
import {MatSortModule} from "@angular/material/sort";
import {NetworkAddComponent} from "./networks/network-add.component";
import {MatFormFieldModule} from "@angular/material/form-field";
import {MatInputModule} from "@angular/material/input";
import {MatGridListModule} from "@angular/material/grid-list";
import {MatDialogModule} from "@angular/material/dialog";
import {NetworkComponent} from "./networks/network.component";
import {NetworkAgentsComponent} from "./networks/network-agents.component";
import {
    EnrollmentRequestApproveDialog,
    NetworkEnrollmentRequestsComponent
} from "./networks/network-enrollment-requests.component";
import {AlertsComponent} from "./alert/alerts.component";
import {MAT_SNACK_BAR_DEFAULT_OPTIONS, MatSnackBarModule} from "@angular/material/snack-bar";

@NgModule({
    declarations: [
        AppComponent,
        AlertsComponent,
        EnrollmentRequestApproveDialog,
        NetworksComponent,
        NetworkComponent,
        NetworkAddComponent,
        NetworkAgentsComponent,
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
        MatSnackBarModule
    ],
    providers: [
        {provide: MAT_SNACK_BAR_DEFAULT_OPTIONS, useValue: {duration: 2500}}
    ],
    bootstrap: [AppComponent]
})
export class AppModule {
}
