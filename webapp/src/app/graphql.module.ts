import {inject, NgModule} from '@angular/core';
import {provideApollo} from 'apollo-angular';
import {ApolloClient, InMemoryCache, ServerError} from "@apollo/client";
import {HttpLink} from 'apollo-angular/http';
import {ErrorLink} from "@apollo/client/link/error";
import {provideHttpClient} from "@angular/common/http";

const uri = '/graphql'; // <-- add the URL of the GraphQL server here
export function createApollo(httpLink: HttpLink): ApolloClient.Options {
    const http = httpLink.create({uri: uri});
    const error = new ErrorLink(({error}) => {
        if (error instanceof ServerError &&
            error.statusCode === 401) {
            window.location.href = '/login';
        }
    });

    const link = error.concat(http);

    return {
        link: link,
        cache: new InMemoryCache({
            typePolicies: {
                Agent: {
                    keyFields: ["fingerprint"]
                },
                EnrollmentRequest: {
                    keyFields: ["fingerprint"]
                },
                Network: {
                    keyFields: ["name"]
                }
            }
        }),
    };
}

@NgModule({
    imports: [],
    providers: [
        provideHttpClient(),
        provideApollo(() => {

            const httpLink = inject(HttpLink);

            return createApollo(httpLink);
        })
    ],
})
export class GraphQLModule {
}
