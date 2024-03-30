import {Suspense} from "react";
import {ErrorBoundary} from "react-error-boundary";
import Card from "../components/card";
import Error from "../components/error";
import Stat from "../components/stat";
import * as Icons from "@heroicons/react/outline";
import Spinner from "../components/spinner";

export default function Home() {
    return (
        <>
            <h3 className="mt-2 mb-5 text-lg font-medium leading-6 text-slate-900">
                Your accounts
            </h3>

            <div className="grid grid-cols-1 gap-5">
                <ErrorBoundary fallback={<Error>Could not fetch data.</Error>}>
                    <Suspense fallback={<Spinner/>}>
                        <Card>
                            <Stat
                                label="Checking"
                                endpoint="/api/checking"
                                Icon={Icons.CashIcon}
                            />
                        </Card>
                    </Suspense>
                </ErrorBoundary>

                <ErrorBoundary fallback={<Error>Could not fetch data.</Error>}>
                    <Suspense fallback={<Spinner/>}>
                        <Card>
                            <Stat
                                label="Savings"
                                endpoint="/api/savings"
                                Icon={Icons.CurrencyDollarIcon}
                            />
                        </Card>
                    </Suspense>
                </ErrorBoundary>

                <ErrorBoundary fallback={<Error>Could not fetch data.</Error>}>
                    <Suspense fallback={<Spinner/>}>
                        <Card>
                            <Stat
                                label="Credit Card"
                                endpoint="/api/credit"
                                Icon={Icons.CreditCardIcon}
                            />
                        </Card>
                    </Suspense>
                </ErrorBoundary>
            </div>

        </>
    )
        ;
}
