import pandas as pd
from decimal import Decimal, ROUND_HALF_UP
from dataclasses import dataclass

@dataclass
class Transfer:
    debetor: str
    spender: str
    value: Decimal

class SplitwiseProcessor:
    def __init__(self, df: pd.DataFrame):
        self.df = df
        self.df['value'] = self.df['value'].apply(lambda x: Decimal(str(x)))

    def calculate(self) -> list[Transfer]:
        spent = self.df.groupby('spender')['value'].sum()
        owed = self.df.groupby('debetor')['value'].sum()
        balances = spent.add(-owed, fill_value=Decimal('0')).sort_values()

        debtors = [[n, v] for n, v in balances.items() if v < 0]
        creditors = [[n, v] for n, v in balances.items() if v > 0][::-1]

        results = []
        i, j = 0, 0
        while i < len(debtors) and j < len(creditors):
            amount = min(-debtors[i][1], creditors[j][1]).quantize(Decimal('0.01'), ROUND_HALF_UP)
            
            if amount > 0:
                results.append(Transfer(debtors[i][0], creditors[j][0], amount))

            debtors[i][1] += amount
            creditors[j][1] -= amount

            if debtors[i][1] == 0: i += 1
            if creditors[j][1] == 0: j += 1
        return results